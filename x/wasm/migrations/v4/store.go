package v4

import (
	"context"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
)

// StoreContractInfoFn stores contract info
type StoreContractInfoFn func(ctx context.Context, contractAddress sdk.AccAddress, contractInfo *types.ContractInfo)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	storeContractInfoFn StoreContractInfoFn
}

// NewMigrator returns a new Migrator.
func NewMigrator(fn StoreContractInfoFn) Migrator {
	return Migrator{storeContractInfoFn: fn}
}

// Migrate4to5 migrates from version 4 to 5.
// This migration re-serializes all ContractInfo records to ensure compatibility
// with the schema that includes the ibc2_port_id field added in wasmd v0.61.5.
//
// Background:
// - wasmd v0.61.4: ContractInfo had 7 fields (no ibc2_port_id)
// - wasmd v0.61.5: ContractInfo added field 8 (ibc2_port_id)
// - Contracts stored with v0.61.4 cause "proto: illegal wireType 7" panics
//
// This migration:
// 1. Iterates through all stored ContractInfo records
// 2. Unmarshals them (works with old schema due to protobuf compatibility)
// 3. Re-marshals and stores them with the current schema
// 4. Ensures ibc2_port_id is properly initialized (empty string if not set)
func (m Migrator) Migrate4to5(ctx sdk.Context, storeService corestoretypes.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(store), types.ContractKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var contractInfo types.ContractInfo

		// Unmarshal existing contract info
		// This will work even if the data was stored with the old schema (v0.61.4)
		// because protobuf handles missing fields gracefully
		if err := cdc.Unmarshal(iter.Value(), &contractInfo); err != nil {
			// If unmarshal fails, skip this contract
			// This shouldn't happen in normal operation
			continue
		}

		// Re-store the contract info with the current schema
		// This ensures the ibc2_port_id field is properly initialized
		// (as empty string if it wasn't set)
		contractAddress := sdk.AccAddress(iter.Key())
		m.storeContractInfoFn(ctx, contractAddress, &contractInfo)
	}

	return nil
}
