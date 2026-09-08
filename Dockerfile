# syntax=docker/dockerfile:1
# docker build . -t cosmwasm/wasmd:latest
# docker run --rm -it cosmwasm/wasmd:latest /bin/sh

FROM golang:1.25-alpine AS go-builder

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/
# Needs --secret id=gitconfig,src=<gitconfig redirecting wasmvm to the private fork>
ENV GOPRIVATE=github.com/CosmWasm/wasmvm
RUN --mount=type=secret,id=gitconfig,target=/root/.gitconfig go mod download

# The private fork ships the muslc static libs as .xz inside the Go module.
RUN ARCH="$(uname -m)" \
  && unxz -c "$(go list -m -f '{{.Dir}}' github.com/CosmWasm/wasmvm/v3)/internal/api/libwasmvm_muslc.$ARCH.a.xz" > "/lib/libwasmvm_muslc.$ARCH.a"

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file /code/build/wasmd | grep "statically linked")

# --------------------------------------------------------
FROM alpine:3.18

COPY --from=go-builder /code/build/wasmd /usr/bin/wasmd

COPY docker/* /opt/
RUN chmod +x /opt/*.sh

WORKDIR /opt

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657

CMD ["/usr/bin/wasmd", "version"]
