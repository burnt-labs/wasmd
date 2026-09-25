# Security Policy

This repository is XION's fork of
[`CosmWasm/wasmd`](https://github.com/CosmWasm/wasmd), and XION mainnet builds
against it. The repository is an asset in the
[Blockchain / DLT bug bounty program](https://github.com/burnt-labs/bug-bounty/blob/main/programs/blockchain.md)
for **Burnt Labs' patches only**.

This file summarizes repository-specific terms. The published
[`burnt-labs/bug-bounty`](https://github.com/burnt-labs/bug-bounty) program is
canonical; where the documents differ, the program terms govern.

## Reporting a Vulnerability

**Do not open a public GitHub issue for a security vulnerability.** Report it
through **Security → Report a vulnerability** on this repository, or email
[security@burnt.com](mailto:security@burnt.com).

We acknowledge receipt within **5 business days** and provide a triage decision
within **14 days**. Active exploitation, or confirmed attacker awareness of an
unpatched vulnerability, escalates the issue to Critical **response handling**
— prioritization, coordination, and disclosure timing — regardless of its
original classification. That escalation does not change the finding's
severity assessment or reward eligibility.

## Fork Scope

Only the delta between this fork and its upstream base is in scope. XION fork
tags use the form `<upstream-tag>-xion.N`; remove the `-xion.N` suffix to
identify the upstream base and diff against that tag. A finding that reproduces
on the unmodified upstream base belongs to CosmWasm and is not eligible under
this program, regardless of its impact on XION.

Scope applies to the fork version in the current XION mainnet release. Findings
affecting only deprecated versions, or already remediated in the currently
deployed release, are not eligible. Verify the version used by the current
mainnet release before submitting.

## Proof of Concept

An end-to-end proof of concept is required. Unit tests or keeper harnesses that
bypass transaction encoding, routing, the ante handler chain, or block execution
do not demonstrate on-chain exploitability on their own.

Run the proof of concept against a locally running XION node configured with
mainnet parameters and execute the attack through standard transaction
broadcast. Broadcast acceptance alone is not sufficient: show inclusion in a
block, the successful execution result, and the resulting state change or
security impact.

## Permissioned Chain Policy

XION mainnet operates with `code_upload_access: Nobody`. Uploading new contract
code requires governance approval. An attack that depends on uploading
attacker-controlled contract code to mainnet is out of scope. A finding that is
exploitable through code already approved for mainnet is not excluded by this
rule.

## Privileged Actor Policy

Findings are classified at **Medium at most** when the attack must begin with
control of governance, a module authority, validator or operator credentials,
or another privileged role — or requires that holder to cooperate — and the
demonstrated action is already within that role's intended authority.

The cap does not apply when a flaw lets an attacker who starts without that
privilege obtain it or bypass its authorization check, or lets a legitimately
held limited role perform actions outside its intended permissions. Those
findings are assessed by demonstrated impact. This policy does not authorize
researchers to acquire or exercise production privileges they do not
legitimately control, or to test with production privileges they do control.

## Rewards and Severity

Only **High** and **Critical** findings are reward eligible. The canonical
program defines severity, exclusions, KYC, duplicate handling, and all other
reward terms.

## Responsible Disclosure and Safe Harbor

Do not test against XION mainnet or other production systems. Use a local
environment or infrastructure you control, do not access or disclose user data,
do not disrupt services, and keep the finding private until disclosure is
coordinated.

Naming this repository as an asset establishes eligibility, not permission to
test a production deployment. Good-faith research within the authorized local
or researcher-controlled environments is covered by the canonical program's
safe harbor. Reporting a vulnerability encountered incidentally is always
welcome.
