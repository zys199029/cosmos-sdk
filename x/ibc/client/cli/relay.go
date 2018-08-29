package cli

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/log"
)

// flags
const (
	FlagStoreName     = "store-name"
	FlagSrcChainID    = "src-chain-id"
	FlagSrcChainNode  = "src-chain-node"
	FlagDestChainID   = "dest-chain-id"
	FlagDestChainNode = "dest-chain-node"
)

type relayCommander struct {
	cdc       *wire.Codec
	address   sdk.AccAddress
	decoder   auth.AccountDecoder
	mainStore string
	accStore  string

	storeName     string
	srcChainID    string
	srcChainNode  string
	destChainID   string
	destChainNode string

	logger log.Logger
}

// IBCRelayCmd implements the IBC relay command.
func IBCRelayCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := relayCommander{
		cdc:       cdc,
		decoder:   authcmd.GetAccountDecoder(cdc),
		mainStore: "main",
		accStore:  "acc",

		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}

	cmd := &cobra.Command{
		Use: "relay",
		Run: cmdr.runIBCRelay,
	}

	cmd.Flags().String(FlagStoreName, "", "Store name to relay packets")
	cmd.Flags().String(FlagSrcChainID, "", "Chain ID for ibc node to check outgoing packets")
	cmd.Flags().String(FlagSrcChainNode, "tcp://localhost:26657", "<host>:<port> to tendermint rpc interface for this chain")
	cmd.Flags().String(FlagDestChainID, "", "Chain ID for ibc node to broadcast incoming packets")
	cmd.Flags().String(FlagDestChainNode, "tcp://localhost:36657", "<host>:<port> to tendermint rpc interface for this chain")

	cmd.MarkFlagRequired(FlagStoreName)
	cmd.MarkFlagRequired(FlagSrcChainID)
	cmd.MarkFlagRequired(FlagSrcChainNode)
	cmd.MarkFlagRequired(FlagDestChainID)
	cmd.MarkFlagRequired(FlagDestChainNode)

	viper.BindPFlag(FlagStoreName, cmd.Flags().Lookup(FlagStoreName))
	viper.BindPFlag(FlagSrcChainID, cmd.Flags().Lookup(FlagSrcChainID))
	viper.BindPFlag(FlagSrcChainNode, cmd.Flags().Lookup(FlagSrcChainNode))
	viper.BindPFlag(FlagDestChainID, cmd.Flags().Lookup(FlagDestChainID))
	viper.BindPFlag(FlagDestChainNode, cmd.Flags().Lookup(FlagDestChainNode))

	return cmd
}

// nolint: unparam
func (c relayCommander) runIBCRelay(cmd *cobra.Command, args []string) {
	storeName := viper.GetString(FlagStoreName)
	srcChainID := viper.GetString(FlagSrcChainID)
	srcChainNode := viper.GetString(FlagSrcChainNode)
	destChainID := viper.GetString(FlagDestChainID)
	destChainNode := viper.GetString(FlagDestChainNode)
	address, err := context.NewCLIContext().GetFromAddress()
	if err != nil {
		panic(err)
	}

	c.storeName = storeName
	c.srcChainID = srcChainID
	c.srcChainNode = srcChainNode
	c.destChainID = destChainID
	c.destChainNode = destChainNode
	c.address = address

	// TODO: use proper config

	c.loop()
}

func (c relayCommander) query(ctx context.CLIContext, path string, params interface{}, ptr interface{}) {
	bz, err := c.cdc.MarshalJSON(params)
	if err != nil {
		panic(err)
	}

	res, err := ctx.QueryWithData("/custom/"+c.storeName+"/ibc/"+path, bz)
	if err != nil {
		panic(err)
	}

	err = c.cdc.UnmarshalJSON(res, ptr)
	if err != nil {
		panic(err)
	}
}

func (c relayCommander) ingressSequence(ctx context.CLIContext) (res uint64) {
	c.query(ctx, "ingress-sequence", ibc.QueryIngressSequenceParams{c.srcChainID, byte(ibc.PacketType)}, &res)
	return
}

func (c relayCommander) egressSequence(ctx context.CLIContext) (res uint64) {
	c.query(ctx, "egress-sequence", ibc.QueryEgressSequenceParams{c.destChainID, byte(ibc.PacketType)}, &res)
	return
}

func (c relayCommander) egressDatagram(ctx context.CLIContext, index uint64) (res ibc.Datagram) {
	c.query(ctx, "egress-datagram", ibc.QueryEgressDatagramParams{c.destChainID, byte(ibc.PacketType), index}, &res)
	return
}

// This is nolinted as someone is in the process of refactoring this to remove the goto
// nolint: gocyclo
func (c relayCommander) loop() {
	authctx := authctx.NewTxContextFromCLI().WithCodec(c.cdc)
	ctx := context.NewCLIContext()
	for {
		time.Sleep(5 * time.Second)

		processed := c.ingressSequence(ctx)
		sequence := c.egressSequence(ctx)
		if sequence <= processed {
			continue
		}

		c.logger.Info("Detected IBC packet", "number", sequence-1)

		var data ibc.Datagram
		for i := processed; i < sequence; i++ {
			// TODO: add proof
			msg := ibc.MsgReceive{Datagram: data, Relayer: c.address}
			err := utils.SendTx(authctx, ctx, []sdk.Msg{msg})
			if err != nil {
				panic(err)
			}

			c.logger.Info("Relayed IBC packet", "number", i)
		}
	}
}
