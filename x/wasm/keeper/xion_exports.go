package keeper

import (
	"context"

	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
// position starts from 0
func (k Keeper) appendToContractHistoryGenesis(ctx context.Context, contractAddr sdk.AccAddress, newEntries ...types.ContractCodeHistoryEntry) error {
	store := k.storeService.OpenKVStore(ctx)
	for pos, e := range newEntries {
		key := types.GetContractCodeHistoryElementKey(contractAddr, uint64(pos))
		if err := store.Set(key, k.cdc.MustMarshal(&e)); err != nil {
			return err
		}
	}
	return nil
}
