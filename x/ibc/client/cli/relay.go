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

type commander struct {
	cdc     *wire.Codec
	decoder auth.AccountDecoder

	address       sdk.AccAddress
	storeName     string
	srcChainID    string
	srcChainNode  string
	destChainID   string
	destChainNode string

	ich chan uint64
	ech chan uint64
	dch chan ibc.Datagram

	logger log.Logger
}

// IBCRelayCmd implements the IBC relay command.
func IBCRelayCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		cdc:     cdc,
		decoder: authcmd.GetAccountDecoder(cdc),

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

func (c commander) runIBCRelay(cmd *cobra.Command, args []string) {
	c.storeName = viper.GetString(FlagStoreName)
	c.srcChainID = viper.GetString(FlagSrcChainID)
	c.srcChainNode = viper.GetString(FlagSrcChainNode)
	c.destChainID = viper.GetString(FlagDestChainID)
	c.destChainNode = viper.GetString(FlagDestChainNode)

	c.ich = make(chan uint64)
	c.ech = make(chan uint64)
	c.dch = make(chan ibc.Datagram)

	ctx := context.NewCLIContext()

	var err error
	c.address, err = ctx.GetFromAddress()
	if err != nil {
		panic(err)
	}

	c.loop(func() { c.queryEgressSequence(ctx) })
	c.loop(func() { c.queryIngressSequence(ctx) })
	c.loop(func() { c.queryDatagram(ctx) })
	c.loop(func() { c.txDatagram(ctx) })

}

func (c commander) queryEgressSequence(ctx context.CLIContext) {
	for {
		seq := c.egressSequence(ctx, ibc.PacketType)
		c.ech <- seq
		c.logger.Info("Detected new packet", "number", seq)
	}
}

func (c commander) queryIngressSequence(ctx context.CLIContext) {
	for {
		seq := c.ingressSequence(ctx, ibc.PacketType)
		c.ich <- seq
		c.logger.Info("Detected processed packet", "number", seq)
	}
}

func (c commander) queryDatagram(ctx context.CLIContext) {
	ingseq := <-c.ich
	egseq := <-c.ech

	for {
		select {
		case newseq := <-c.ich:
			ingseq = newseq
		case newseq := <-c.ech:
			egseq = newseq
		default:
			if egseq <= ingseq {
				time.Sleep(5 * time.Second)
				continue
			}
			for seq := ingseq; seq < egseq; seq++ {
				data := c.egressDatagram(ctx, ingseq, ibc.PacketType)
				c.logger.Info("Retrieved packet", "number", ingseq)
				c.dch <- data
			}
			ingseq = egseq
			time.Sleep(1 * time.Second)
		}
	}
}

func (c commander) txDatagram(ctx context.CLIContext) {
	for {
		data := <-c.dch
		msg := ibc.MsgReceive{Datagram: data, Relayer: c.address}
		err := c.tx(ctx, msg)
		c.logger.Info("Submitted packet")
		if err != nil {
			c.logger.Info("Error transacting packet, skipping")
		}
	}
}

func (c commander) tx(ctx context.CLIContext, msg sdk.Msg) error {
	authctx := authctx.NewTxContextFromCLI().WithCodec(c.cdc)
	return utils.SendTx(authctx, ctx, []sdk.Msg{msg})
}

func (c commander) ingressSequence(ctx context.CLIContext, ty ibc.DatagramType) (res uint64) {
	c.query(ctx, "ingress-sequence", ibc.QueryIngressSequenceParams{c.srcChainID, byte(ty)}, &res)
	return
}

func (c commander) egressSequence(ctx context.CLIContext, ty ibc.DatagramType) (res uint64) {
	c.query(ctx, "egress-sequence", ibc.QueryEgressSequenceParams{c.destChainID, byte(ty)}, &res)
	return
}

func (c commander) egressDatagram(ctx context.CLIContext, index uint64, ty ibc.DatagramType) (res ibc.Datagram) {
	c.query(ctx, "egress-datagram", ibc.QueryEgressDatagramParams{c.destChainID, byte(ty), index}, &res)
	return
}

func (c commander) loop(f func()) {
	go func() {
		for {
			func() {
				defer func() {
					if err := recover(); err != nil {
						c.logger.Info("Panic!", "msg", err)
					}
				}()
				f()
			}()
		}
	}()
}

func (c commander) query(ctx context.CLIContext, path string, params interface{}, ptr interface{}) {
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
