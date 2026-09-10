package keeper

import (
	"os"
	"testing"

	wasmvm "github.com/CosmWasm/wasmvm/v3"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v3/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	"github.com/CosmWasm/wasmd/x/wasm/keeper/wasmtesting"
	"github.com/CosmWasm/wasmd/x/wasm/types"
)

// A wasmvm ContractResult always sets exactly one of ok/error. A result with
// neither set is malformed output — from a buggy wasmvm or a contract whose
// return value round-trips to `{}`. Every entrypoint below used to dereference
// res.Ok straight after the res.Err check, so such a result nil-panicked the
// node. Inside a tx that panic is recovered by baseapp, but any caller running
// outside runTx — a Begin/End/PreBlocker executing a contract — halts the
// chain. These tests pin the guards: each call must return ErrVMError, and must
// not panic.

// nilOkResult is the malformed value: neither Ok nor Err populated.
func nilOkResult() *wasmvmtypes.ContractResult { return &wasmvmtypes.ContractResult{} }

func okResult() *wasmvmtypes.ContractResult {
	return &wasmvmtypes.ContractResult{Ok: &wasmvmtypes.Response{}}
}

// newNilOkKeeper returns a keeper whose engine answers every entrypoint well
// enough to instantiate a contract, plus the instantiated contract address. The
// caller then swaps in the malformed response for the entrypoint under test.
func newNilOkKeeper(t *testing.T) (sdk.Context, *wasmtesting.MockWasmEngine, *TestKeepers, sdk.AccAddress) {
	t.Helper()

	mock := &wasmtesting.MockWasmEngine{}
	wasmtesting.MakeInstantiable(mock)
	mock.MigrateWithInfoFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ []byte, _ wasmvmtypes.MigrateInfo, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return okResult(), 1, nil
	}
	mock.MigrateFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ []byte, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return okResult(), 1, nil
	}

	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities, WithWasmEngine(mock))
	creator := DeterministicAccountAddress(t, 1)
	keepers.Faucet.Fund(ctx, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))...)

	codeID, _, err := keepers.ContractKeeper.Create(ctx, creator, hackatomWasm, nil)
	require.NoError(t, err)

	// admin is set so the migrate test reaches the entrypoint rather than
	// stopping at the authorization check
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx, codeID, creator, creator, []byte(`{}`), "test", nil)
	require.NoError(t, err)

	return ctx, mock, &keepers, contractAddr
}

func TestInstantiateRejectsNilOkResponse(t *testing.T) {
	mock := &wasmtesting.MockWasmEngine{}
	wasmtesting.MakeInstantiable(mock)
	mock.InstantiateFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ wasmvmtypes.MessageInfo, _ []byte, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}

	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities, WithWasmEngine(mock))
	creator := DeterministicAccountAddress(t, 1)
	keepers.Faucet.Fund(ctx, creator, sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))...)

	codeID, _, err := keepers.ContractKeeper.Create(ctx, creator, hackatomWasm, nil)
	require.NoError(t, err)

	require.NotPanics(t, func() {
		_, _, err = keepers.ContractKeeper.Instantiate(ctx, codeID, creator, nil, []byte(`{}`), "test", nil)
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

func TestExecuteRejectsNilOkResponse(t *testing.T) {
	ctx, mock, keepers, contractAddr := newNilOkKeeper(t)
	mock.ExecuteFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ wasmvmtypes.MessageInfo, _ []byte, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}

	var err error
	require.NotPanics(t, func() {
		_, err = keepers.ContractKeeper.Execute(ctx, contractAddr, DeterministicAccountAddress(t, 1), []byte(`{}`), nil)
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

func TestSudoRejectsNilOkResponse(t *testing.T) {
	ctx, mock, keepers, contractAddr := newNilOkKeeper(t)
	mock.SudoFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ []byte, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}

	var err error
	require.NotPanics(t, func() {
		_, err = keepers.ContractKeeper.Sudo(ctx, contractAddr, []byte(`{}`))
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

func TestReplyRejectsNilOkResponse(t *testing.T) {
	ctx, mock, keepers, contractAddr := newNilOkKeeper(t)
	mock.ReplyFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ wasmvmtypes.Reply, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}

	var err error
	require.NotPanics(t, func() {
		_, err = keepers.WasmKeeper.reply(ctx, contractAddr, wasmvmtypes.Reply{})
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

func TestMigrateRejectsNilOkResponse(t *testing.T) {
	ctx, mock, keepers, contractAddr := newNilOkKeeper(t)
	creator := DeterministicAccountAddress(t, 1)

	burnerWasm, err := os.ReadFile("./testdata/burner.wasm")
	require.NoError(t, err)
	newCodeID, _, err := keepers.ContractKeeper.Create(ctx, creator, burnerWasm, nil)
	require.NoError(t, err)

	mock.MigrateWithInfoFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ []byte, _ wasmvmtypes.MigrateInfo, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}
	mock.MigrateFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ []byte, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.ContractResult, uint64, error) {
		return nilOkResult(), 1, nil
	}

	require.NotPanics(t, func() {
		_, err = keepers.ContractKeeper.Migrate(ctx, contractAddr, creator, newCodeID, []byte(`{}`))
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

// handleIBCBasicContractResponse is where the nine IBC call sites that pass
// res.Ok straight through converge, so one guard covers them all.
func TestHandleIBCBasicContractResponseRejectsNil(t *testing.T) {
	ctx, _, keepers, contractAddr := newNilOkKeeper(t)

	var err error
	require.NotPanics(t, func() {
		err = keepers.WasmKeeper.handleIBCBasicContractResponse(ctx, contractAddr, "port", nil)
	})
	require.ErrorIs(t, err, types.ErrVMError)
}

// nilOkReceiveResult is the malformed value for the packet-receive entrypoints:
// an IBCReceiveResult with neither Ok nor Err populated.
func nilOkReceiveResult() *wasmvmtypes.IBCReceiveResult {
	return &wasmvmtypes.IBCReceiveResult{}
}

// newIBCContract seeds an IBC-capable contract, which the packet-receive
// entrypoints require.
func newIBCContract(t *testing.T) (sdk.Context, *wasmtesting.MockWasmEngine, *TestKeepers, sdk.AccAddress) {
	t.Helper()

	var m wasmtesting.MockWasmEngine
	wasmtesting.MakeIBCInstantiable(&m)
	ctx, keepers := CreateTestInput(t, false, AvailableCapabilities)
	example := SeedNewContractInstance(t, ctx, keepers, &m)

	return ctx, &m, &keepers, example.Contract
}

// OnRecvPacket has its own guard rather than reaching the shared basic-response
// helper, and returns (nil, error) rather than a failure ack, so it needs its
// own case.
func TestOnRecvPacketRejectsNilOkResponse(t *testing.T) {
	ctx, m, keepers, contractAddr := newIBCContract(t)
	m.IBCPacketReceiveFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ wasmvmtypes.IBCPacketReceiveMsg, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.IBCReceiveResult, uint64, error) {
		return nilOkReceiveResult(), 1, nil
	}

	var (
		ack ibcexported.Acknowledgement
		err error
	)
	require.NotPanics(t, func() {
		ack, err = keepers.WasmKeeper.OnRecvPacket(ctx, contractAddr,
			wasmvmtypes.IBCPacketReceiveMsg{Packet: wasmvmtypes.IBCPacket{Data: []byte("my data")}})
	})
	require.ErrorIs(t, err, types.ErrVMError)
	require.Nil(t, ack)
}

// OnRecvIBC2Packet returns a failure RecvPacketResult instead of an error, so
// assert that shape specifically rather than reusing the error-returning case.
func TestOnRecvIBC2PacketRejectsNilOkResponse(t *testing.T) {
	ctx, m, keepers, contractAddr := newIBCContract(t)
	m.IBC2PacketReceiveFn = func(_ wasmvm.Checksum, _ wasmvmtypes.Env, _ wasmvmtypes.IBC2PacketReceiveMsg, _ wasmvm.KVStore, _ wasmvm.GoAPI, _ wasmvm.Querier, _ wasmvm.GasMeter, _ uint64, _ wasmvmtypes.UFraction) (*wasmvmtypes.IBCReceiveResult, uint64, error) {
		return nilOkReceiveResult(), 1, nil
	}

	var res channeltypesv2.RecvPacketResult
	require.NotPanics(t, func() {
		res = keepers.WasmKeeper.OnRecvIBC2Packet(ctx, contractAddr,
			wasmvmtypes.IBC2PacketReceiveMsg{Payload: wasmvmtypes.IBC2Payload{Value: []byte("my data")}})
	})
	require.Equal(t, channeltypesv2.PacketStatus_Failure, res.Status)
	require.Contains(t, string(res.Acknowledgement), "nil ok response")
}
