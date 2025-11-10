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
// CRITICAL: This function MUST be deterministic across all validators.
// We always use proto.Unmarshal (not cdc.Unmarshal) to ensure consistent behavior.
// The proto.Unmarshal is more lenient and handles missing fields gracefully.
func SafeUnmarshalContractInfo(cdc codec.BinaryCodec, bz []byte, contractInfo *types.ContractInfo) error {
	// ALWAYS use proto.Unmarshal for deterministic behavior
	// This works with both old schema (v0.61.4) and new schema (v0.61.5+)
	// Missing fields are initialized to zero values (empty string for ibc2_port_id)
	if err := proto.Unmarshal(bz, contractInfo); err != nil {
		return fmt.Errorf("failed to unmarshal ContractInfo: %w", err)
	}
	return nil
}

// SafeUnmarshalCodeInfo attempts to unmarshal CodeInfo with support for legacy formats.
// This provides symmetric handling for CodeInfo similar to ContractInfo, ensuring
// that any schema changes to CodeInfo are also handled gracefully.
//
// CRITICAL: This function MUST be deterministic across all validators.
// We always use proto.Unmarshal (not cdc.Unmarshal) to ensure consistent behavior.
func SafeUnmarshalCodeInfo(cdc codec.BinaryCodec, bz []byte, codeInfo *types.CodeInfo) error {
	// ALWAYS use proto.Unmarshal for deterministic behavior
	if err := proto.Unmarshal(bz, codeInfo); err != nil {
		return fmt.Errorf("failed to unmarshal CodeInfo: %w", err)
	}
	return nil
}
