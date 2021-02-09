package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/band-consumer/x/consuming/types"
	bandoracle "github.com/bandprotocol/chain/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) Result(c context.Context, req *types.QueryResultRequest) (*types.QueryResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bz, err := q.GetResult(ctx, bandoracle.RequestID(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &types.QueryResultResponse{Result: bz}, nil
}

func (q Keeper) LatestRequestID(c context.Context, req *types.QueryLatestRequestIDRequest) (*types.QueryLatestRequestIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := q.GetLatestRequestID(ctx)
	return &types.QueryLatestRequestIDResponse{RequestId: id}, nil
}
