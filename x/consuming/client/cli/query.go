package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/band-consumer/x/consuming/types"
)

// GetQueryCmd returns
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the consuming module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(
		GetCmdReadResult(),
		GetCmdLatestRequest(),
	)

	return queryCmd
}

// GetCmdReadResult queries request by reqID
func GetCmdReadResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "result",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			r, err := queryClient.Result(context.Background(), &types.QueryResultRequest{RequestId: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdLatestRequest queries latest request
func GetCmdLatestRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use: "latest-request",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			r, err := queryClient.LatestRequestID(context.Background(), &types.QueryLatestRequestIDRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(r)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
