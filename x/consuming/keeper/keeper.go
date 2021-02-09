package keeper

import (
	"fmt"

	oracletypes "github.com/bandprotocol/chain/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/bandprotocol/band-consumer/x/consuming/types"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	scopedKeeper  capabilitykeeper.ScopedKeeper
}

// NewKeeper creates a new oracle Keeper instance.
func NewKeeper(cdc codec.BinaryMarshaler, key sdk.StoreKey, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper, scopedKeeper capabilitykeeper.ScopedKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

func (k Keeper) SetResult(ctx sdk.Context, requestID oracletypes.RequestID, result []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ResultStoreKey(requestID), result)
}

func (k Keeper) GetResult(ctx sdk.Context, requestID oracletypes.RequestID) ([]byte, error) {
	if !k.HasResult(ctx, requestID) {
		return nil, sdkerrors.Wrapf(types.ErrItemNotFound,
			"GetResult: Result for request ID %d is not available.", requestID,
		)
	}
	store := ctx.KVStore(k.storeKey)
	return store.Get(types.ResultStoreKey(requestID)), nil
}

func (k Keeper) HasResult(ctx sdk.Context, requestID oracletypes.RequestID) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ResultStoreKey(requestID))
}

func (k Keeper) GetLatestRequestID(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.LatestRequestIDKey)
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

func (k Keeper) SetLatestRequestID(ctx sdk.Context, id oracletypes.RequestID) {
	ctx.KVStore(k.storeKey).Set(types.LatestRequestIDKey, k.cdc.MustMarshalBinaryLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	fmt.Println(portID)
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}
