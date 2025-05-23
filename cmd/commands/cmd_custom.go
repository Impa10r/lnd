package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/urfave/cli"
)

var sendCustomCommand = cli.Command{
	Name:     "sendcustom",
	Category: "Peers",
	Usage:    "Send a custom p2p wire message to a peer",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "peer",
		},
		cli.Uint64Flag{
			Name: "type",
		},
		cli.StringFlag{
			Name: "data",
		},
	},
	Action: actionDecorator(sendCustom),
}

func sendCustom(ctx *cli.Context) error {
	ctxc := getContext()
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	peer, err := hex.DecodeString(ctx.String("peer"))
	if err != nil {
		return err
	}

	msgType := ctx.Uint64("type")

	data, err := hex.DecodeString(ctx.String("data"))
	if err != nil {
		return err
	}

	resp, err := client.SendCustomMessage(
		ctxc, &lnrpc.SendCustomMessageRequest{
			Peer: peer,
			Type: uint32(msgType),
			Data: data,
		},
	)

	printRespJSON(resp)

	return err
}

var subscribeCustomCommand = cli.Command{
	Name:     "subscribecustom",
	Category: "Peers",
	Usage: "Subscribe to incoming custom p2p wire messages from all " +
		"peers",
	Action: actionDecorator(subscribeCustom),
}

func subscribeCustom(ctx *cli.Context) error {
	ctxc := getContext()
	client, cleanUp := getClient(ctx)
	defer cleanUp()

	stream, err := client.SubscribeCustomMessages(
		ctxc, &lnrpc.SubscribeCustomMessagesRequest{},
	)
	if err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		fmt.Printf("Received from peer %x: type=%d, data=%x\n",
			msg.Peer, msg.Type, msg.Data)
	}
}
