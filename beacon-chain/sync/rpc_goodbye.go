package sync

import (
	"context"
	"fmt"
	"time"

	libp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	p2ptypes "github.com/prysmaticlabs/prysm/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/shared/mputil"
	"github.com/sirupsen/logrus"
)

var backOffTime = map[types.SSZUint64]time.Duration{
	// Do not dial peers which are from a different/unverifiable
	// network.
	p2ptypes.GoodbyeCodeWrongNetwork:          24 * time.Hour,
	p2ptypes.GoodbyeCodeUnableToVerifyNetwork: 24 * time.Hour,
	// If local peer is banned, we back off for
	// 2 hours to let the remote peer score us
	// back up again.
	p2ptypes.GoodbyeCodeBadScore:       2 * time.Hour,
	p2ptypes.GoodbyeCodeBanned:         2 * time.Hour,
	p2ptypes.GoodbyeCodeClientShutdown: 1 * time.Hour,
	// Wait 5 minutes before dialing a peer who is
	// 'full'
	p2ptypes.GoodbyeCodeTooManyPeers: 5 * time.Minute,
	p2ptypes.GoodbyeCodeGenericError: 2 * time.Minute,
}

// goodbyeRPCHandler reads the incoming goodbye rpc message from the peer.
func (s *Service) goodbyeRPCHandler(_ context.Context, msg interface{}, stream libp2pcore.Stream) error {
	SetRPCStreamDeadlines(stream)

	m, ok := msg.(*types.SSZUint64)
	if !ok {
		return fmt.Errorf("wrong message type for goodbye, got %T, wanted *uint64", msg)
	}
	if err := s.rateLimiter.validateRequest(stream, 1); err != nil {
		return err
	}
	s.rateLimiter.add(stream, 1)
	log := log.WithField("Reason", goodbyeMessage(*m))
	log.WithField("peer", stream.Conn().RemotePeer()).Debug("Peer has sent a goodbye message")
	s.cfg.P2P.Peers().SetNextValidTime(stream.Conn().RemotePeer(), goodByeBackoff(*m))
	// closes all streams with the peer
	return s.cfg.P2P.Disconnect(stream.Conn().RemotePeer())
}

// disconnectBadPeer checks whether peer is considered bad by some scorer, and tries to disconnect
// the peer, if that is the case. Additionally, disconnection reason is obtained from scorer.
func (s *Service) disconnectBadPeer(ctx context.Context, id peer.ID) {
	if !s.cfg.P2P.Peers().IsBad(id) {
		return
	}
	goodbyeCode := p2ptypes.ErrToGoodbyeCode(s.cfg.P2P.Peers().Scorers().ValidationError(id))
	if err := s.sendGoodByeAndDisconnect(ctx, goodbyeCode, id); err != nil {
		log.Debugf("Error when disconnecting with bad peer: %v", err)
	}
}

// A custom goodbye method that is used by our connection handler, in the
// event we receive bad peers.
func (s *Service) sendGoodbye(ctx context.Context, id peer.ID) error {
	return s.sendGoodByeAndDisconnect(ctx, p2ptypes.GoodbyeCodeGenericError, id)
}

func (s *Service) sendGoodByeAndDisconnect(ctx context.Context, code p2ptypes.RPCGoodbyeCode, id peer.ID) error {
	lock := mputil.NewMultilock(id.String())
	lock.Lock()
	defer lock.Unlock()
	// In the event we are already disconnected, exit early from the
	// goodbye method to prevent redundant streams from being created.
	if s.cfg.P2P.Host().Network().Connectedness(id) == network.NotConnected {
		return nil
	}
	if err := s.sendGoodByeMessage(ctx, code, id); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"peer":  id,
		}).Debug("Could not send goodbye message to peer")
	}
	return s.cfg.P2P.Disconnect(id)
}

func (s *Service) sendGoodByeMessage(ctx context.Context, code p2ptypes.RPCGoodbyeCode, id peer.ID) error {
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()

	topic, err := p2p.TopicFromMessage(p2p.GoodbyeMessageName, helpers.SlotToEpoch(s.cfg.Chain.CurrentSlot()))
	if err != nil {
		return err
	}
	stream, err := s.cfg.P2P.Send(ctx, &code, topic, id)
	if err != nil {
		return err
	}
	defer closeStream(stream, log)

	log := log.WithField("Reason", goodbyeMessage(code))
	log.WithField("peer", stream.Conn().RemotePeer()).Debug("Sending Goodbye message to peer")

	// Wait up to the response timeout for the peer to receive the goodbye
	// and close the stream (or disconnect). We usually don't bother waiting
	// around for an EOF, but we're going to close this connection
	// immediately after we say goodbye.
	//
	// NOTE: we don't actually check the response as there's nothing we can
	// do if something fails. We just need to wait for it.
	SetStreamReadDeadline(stream, respTimeout)
	_, _err := stream.Read([]byte{0})
	_ = _err

	return nil
}

func goodbyeMessage(num p2ptypes.RPCGoodbyeCode) string {
	reason, ok := p2ptypes.GoodbyeCodeMessages[num]
	if ok {
		return reason
	}
	return fmt.Sprintf("unknown goodbye value of %d received", num)
}

// determines which backoff time to use depending on the
// goodbye code provided.
func goodByeBackoff(num p2ptypes.RPCGoodbyeCode) time.Time {
	duration, ok := backOffTime[num]
	if !ok {
		return time.Time{}
	}
	return time.Now().Add(duration)
}
