package p2p

import (
	"context"

	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/network/forks"
	pb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var responseCodeSuccess = byte(0x00)

func (c *client) registerHandshakeHandlers() {
	c.registerRPCHandler(p2p.RPCPingTopicV1, c.pingHandler)
	c.registerRPCHandler(p2p.RPCStatusTopicV1, c.statusRPCHandler)
}

// pingHandler reads the incoming ping rpc message from the peer.
func (c *client) pingHandler(_ context.Context, _ interface{}, stream libp2pcore.Stream) error {
	defer closeStream(stream)
	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		return err
	}
	sq := types.SSZUint64(c.MetadataSeq())
	if _, err := c.Encoding().EncodeWithMaxLength(stream, &sq); err != nil {
		return err
	}
	return nil
}

// statusRPCHandler reads the incoming Status RPC from the peer and responds with our version of a status message.
// This handler will disconnect any peer that does not match our fork version.
func (c *client) statusRPCHandler(ctx context.Context, _ interface{}, stream libp2pcore.Stream) error {
	defer closeStream(stream)
	chainHead, err := c.beaconClient.GetChainHead(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	resp, err := c.nodeClient.GetGenesis(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	digest, err := forks.CreateForkDigest(resp.GenesisTime.AsTime(), resp.GenesisValidatorsRoot)
	if err != nil {
		return err
	}
	status := &pb.Status{
		ForkDigest:     digest[:],
		FinalizedRoot:  chainHead.FinalizedBlockRoot,
		FinalizedEpoch: chainHead.FinalizedEpoch,
		HeadRoot:       chainHead.HeadBlockRoot,
		HeadSlot:       chainHead.HeadSlot,
	}

	if _, err := stream.Write([]byte{responseCodeSuccess}); err != nil {
		log.WithError(err).Debug("Could not write to stream")
		return err
	}
	_, err = c.Encoding().EncodeWithMaxLength(stream, status)
	return err
}