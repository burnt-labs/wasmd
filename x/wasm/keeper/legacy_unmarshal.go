package keeper

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// SafeUnmarshalContractInfo attempts to unmarshal ContractInfo with support for legacy formats.
// This handles the case where contracts were stored with wasmd v0.61.4 (without ibc2_port_id field)
// and need to be read by wasmd v0.61.5+ (which includes the ibc2_port_id field).
//
// The protobuf error "proto: illegal wireType 7" occurs when the new schema tries to read
// data stored with an older schema that's missing fields.
//
// This function:
// 1. First attempts standard unmarshaling with the current schema
// 2. If that fails, attempts to unmarshal using raw protobuf (more lenient)
// 3. Returns an error only if both attempts fail
func SafeUnmarshalContractInfo(cdc codec.BinaryCodec, bz []byte, contractInfo *types.ContractInfo) error {
	// Try standard unmarshal first
	err := cdc.Unmarshal(bz, contractInfo)
	if err == nil {
		return nil
	}

	// If standard unmarshal fails, try raw proto unmarshal
	// This is more lenient and will work with data stored using the old schema
	legacyInfo := &types.ContractInfo{}
	if err := proto.Unmarshal(bz, legacyInfo); err != nil {
		return fmt.Errorf("failed to unmarshal ContractInfo with both current and legacy formats: %w", err)
	}

	// Successfully unmarshaled with legacy format, copy to output
	*contractInfo = *legacyInfo
	return nil
}

// SafeUnmarshalCodeInfo attempts to unmarshal CodeInfo with support for legacy formats.
// This provides symmetric handling for CodeInfo similar to ContractInfo, ensuring
// that any schema changes to CodeInfo are also handled gracefully.
//
// While CodeInfo hasn't changed between v0.61.4 and v0.61.5, this function provides:
// 1. Future-proofing against CodeInfo schema changes
// 2. Consistent error handling across all unmarshal operations
// 3. A clear pattern for handling schema migrations
func SafeUnmarshalCodeInfo(cdc codec.BinaryCodec, bz []byte, codeInfo *types.CodeInfo) error {
	// Try standard unmarshal first
	err := cdc.Unmarshal(bz, codeInfo)
	if err == nil {
		return nil
	}

	// If standard unmarshal fails, try raw proto unmarshal
	legacyInfo := &types.CodeInfo{}
	if err := proto.Unmarshal(bz, legacyInfo); err != nil {
		return fmt.Errorf("failed to unmarshal CodeInfo with both current and legacy formats: %w", err)
	}

	// Successfully unmarshaled with legacy format, copy to output
	*codeInfo = *legacyInfo
	return nil
}
