package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/wasm/types"
)

func (k Keeper) ImportCode(ctx context.Context, codeID uint64, codeInfo types.CodeInfo, wasmCode []byte) error {
	return k.importCode(ctx, codeID, codeInfo, wasmCode)
}

func (k Keeper) ImportContract(ctx context.Context, contractAddr sdk.AccAddress, c *types.ContractInfo, state []types.Model, historyEntries []types.ContractCodeHistoryEntry) error {
	return k.importContract(ctx, contractAddr, c, state, historyEntries)
}

func (k Keeper) ImportAutoIncrementID(ctx context.Context, sequenceKey []byte, val uint64) error {
	return k.importAutoIncrementID(ctx, sequenceKey, val)
}

// appendToContractHistoryGenesis is a helper function to append to the contract history for genesis.
// it skips creating iterators and assumes the contract history is already sorted
// Contract history positions start at 1, matching appendToContractHistory.
func (k Keeper) appendToContractHistoryGenesis(ctx context.Context, contractAddr sdk.AccAddress, newEntries ...types.ContractCodeHistoryEntry) error {
	store := k.storeService.OpenKVStore(ctx)
	for pos, e := range newEntries {
		key := types.GetContractCodeHistoryElementKey(contractAddr, uint64(pos)+1)
		if err := store.Set(key, k.cdc.MustMarshal(&e)); err != nil {
			return err
		}
	}
	return nil
}
