package state_test

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	core "github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	p2pType "github.com/prysmaticlabs/prysm/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	prysmv2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
	"github.com/prysmaticlabs/prysm/proto/prysm/v2/wrapper"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/copyutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestExecuteAltairStateTransitionNoVerify_FullProcess(t *testing.T) {
	beaconState, privKeys := testutil.DeterministicGenesisStateAltair(t, 100)

	syncCommittee, err := altair.NextSyncCommittee(beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(syncCommittee))

	eth1Data := &ethpb.Eth1Data{
		DepositCount: 100,
		DepositRoot:  bytesutil.PadTo([]byte{2}, 32),
		BlockHash:    make([]byte, 32),
	}
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch-1))
	e := beaconState.Eth1Data()
	e.DepositCount = 100
	require.NoError(t, beaconState.SetEth1Data(e))
	bh := beaconState.LatestBlockHeader()
	bh.Slot = beaconState.Slot()
	require.NoError(t, beaconState.SetLatestBlockHeader(bh))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data}))

	require.NoError(t, beaconState.SetSlot(beaconState.Slot()+1))
	epoch := helpers.CurrentEpoch(beaconState)
	randaoReveal, err := testutil.RandaoReveal(beaconState, epoch, privKeys)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetSlot(beaconState.Slot()-1))

	nextSlotState, err := core.ProcessSlots(context.Background(), beaconState.Copy(), beaconState.Slot()+1)
	require.NoError(t, err)
	parentRoot, err := nextSlotState.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(nextSlotState)
	require.NoError(t, err)
	block := testutil.NewBeaconBlockAltair()
	block.Block.ProposerIndex = proposerIdx
	block.Block.Slot = beaconState.Slot() + 1
	block.Block.ParentRoot = parentRoot[:]
	block.Block.Body.RandaoReveal = randaoReveal
	block.Block.Body.Eth1Data = eth1Data

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xff
	}
	indices, err := altair.NextSyncCommitteeIndices(beaconState)
	require.NoError(t, err)
	h := copyutil.CopyBeaconBlockHeader(beaconState.LatestBlockHeader())
	prevStateRoot, err := beaconState.HashTreeRoot(context.Background())
	require.NoError(t, err)
	h.StateRoot = prevStateRoot[:]
	pbr, err := h.HashTreeRoot()
	require.NoError(t, err)
	syncSigs := make([]bls.Signature, len(indices))
	for i, indice := range indices {
		b := p2pType.SSZBytes(pbr[:])
		sb, err := helpers.ComputeDomainAndSign(beaconState, helpers.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
		require.NoError(t, err)
		sig, err := bls.SignatureFromBytes(sb)
		require.NoError(t, err)
		syncSigs[i] = sig
	}
	aggregatedSig := bls.AggregateSignatures(syncSigs).Marshal()
	syncAggregate := &prysmv2.SyncAggregate{
		SyncCommitteeBits:      syncBits,
		SyncCommitteeSignature: aggregatedSig,
	}
	block.Block.Body.SyncAggregate = syncAggregate
	wsb, err := wrapper.WrappedAltairSignedBeaconBlock(block)
	require.NoError(t, err)
	stateRoot, err := core.CalculateStateRoot(context.Background(), beaconState, wsb)
	require.NoError(t, err)
	block.Block.StateRoot = stateRoot[:]

	c := beaconState.Copy()
	sig, err := testutil.BlockSignatureAltair(c, block.Block, privKeys)
	require.NoError(t, err)
	block.Signature = sig.Marshal()

	wsb, err = wrapper.WrappedAltairSignedBeaconBlock(block)
	require.NoError(t, err)
	set, _, err := core.ExecuteStateTransitionNoVerifyAnySig(context.Background(), beaconState, wsb)
	require.NoError(t, err)
	verified, err := set.Verify()
	require.NoError(t, err)
	require.Equal(t, true, verified, "Could not verify signature set")
}

func TestExecuteAltairStateTransitionNoVerifySignature_CouldNotVerifyStateRoot(t *testing.T) {
	beaconState, privKeys := testutil.DeterministicGenesisStateAltair(t, 100)

	syncCommittee, err := altair.NextSyncCommittee(beaconState)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetCurrentSyncCommittee(syncCommittee))

	eth1Data := &ethpb.Eth1Data{
		DepositCount: 100,
		DepositRoot:  bytesutil.PadTo([]byte{2}, 32),
		BlockHash:    make([]byte, 32),
	}
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch-1))
	e := beaconState.Eth1Data()
	e.DepositCount = 100
	require.NoError(t, beaconState.SetEth1Data(e))
	bh := beaconState.LatestBlockHeader()
	bh.Slot = beaconState.Slot()
	require.NoError(t, beaconState.SetLatestBlockHeader(bh))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data}))

	require.NoError(t, beaconState.SetSlot(beaconState.Slot()+1))
	epoch := helpers.CurrentEpoch(beaconState)
	randaoReveal, err := testutil.RandaoReveal(beaconState, epoch, privKeys)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetSlot(beaconState.Slot()-1))

	nextSlotState, err := core.ProcessSlots(context.Background(), beaconState.Copy(), beaconState.Slot()+1)
	require.NoError(t, err)
	parentRoot, err := nextSlotState.LatestBlockHeader().HashTreeRoot()
	require.NoError(t, err)
	proposerIdx, err := helpers.BeaconProposerIndex(nextSlotState)
	require.NoError(t, err)
	block := testutil.NewBeaconBlockAltair()
	block.Block.ProposerIndex = proposerIdx
	block.Block.Slot = beaconState.Slot() + 1
	block.Block.ParentRoot = parentRoot[:]
	block.Block.Body.RandaoReveal = randaoReveal
	block.Block.Body.Eth1Data = eth1Data

	syncBits := bitfield.NewBitvector512()
	for i := range syncBits {
		syncBits[i] = 0xff
	}
	indices, err := altair.NextSyncCommitteeIndices(beaconState)
	require.NoError(t, err)
	h := copyutil.CopyBeaconBlockHeader(beaconState.LatestBlockHeader())
	prevStateRoot, err := beaconState.HashTreeRoot(context.Background())
	require.NoError(t, err)
	h.StateRoot = prevStateRoot[:]
	pbr, err := h.HashTreeRoot()
	require.NoError(t, err)
	syncSigs := make([]bls.Signature, len(indices))
	for i, indice := range indices {
		b := p2pType.SSZBytes(pbr[:])
		sb, err := helpers.ComputeDomainAndSign(beaconState, helpers.CurrentEpoch(beaconState), &b, params.BeaconConfig().DomainSyncCommittee, privKeys[indice])
		require.NoError(t, err)
		sig, err := bls.SignatureFromBytes(sb)
		require.NoError(t, err)
		syncSigs[i] = sig
	}
	aggregatedSig := bls.AggregateSignatures(syncSigs).Marshal()
	syncAggregate := &prysmv2.SyncAggregate{
		SyncCommitteeBits:      syncBits,
		SyncCommitteeSignature: aggregatedSig,
	}
	block.Block.Body.SyncAggregate = syncAggregate

	wsb, err := wrapper.WrappedAltairSignedBeaconBlock(block)
	require.NoError(t, err)
	stateRoot, err := core.CalculateStateRoot(context.Background(), beaconState, wsb)
	require.NoError(t, err)
	block.Block.StateRoot = stateRoot[:]

	c := beaconState.Copy()
	sig, err := testutil.BlockSignatureAltair(c, block.Block, privKeys)
	require.NoError(t, err)
	block.Signature = sig.Marshal()

	block.Block.StateRoot = bytesutil.PadTo([]byte{'a'}, 32)
	wsb, err = wrapper.WrappedAltairSignedBeaconBlock(block)
	require.NoError(t, err)
	_, _, err = core.ExecuteStateTransitionNoVerifyAnySig(context.Background(), beaconState, wsb)
	require.ErrorContains(t, "could not validate state root", err)
}

func TestExecuteStateTransitionNoVerifyAnySig_PassesProcessingConditions(t *testing.T) {
	beaconState, block := createFullAltairBlockWithOperations(t)
	wsb, err := wrapper.WrappedAltairSignedBeaconBlock(block)
	require.NoError(t, err)
	set, _, err := core.ExecuteStateTransitionNoVerifyAnySig(context.Background(), beaconState, wsb)
	require.NoError(t, err)
	// Test Signature set verifies.
	verified, err := set.Verify()
	require.NoError(t, err)
	require.Equal(t, true, verified, "Could not verify signature set")
}

func createFullAltairBlockWithOperations(t *testing.T) (state.BeaconStateAltair,
	*prysmv2.SignedBeaconBlock) {
	beaconState, privKeys := testutil.DeterministicGenesisStateAltair(t, 32)
	sCom, err := altair.NextSyncCommittee(beaconState)
	assert.NoError(t, err)
	assert.NoError(t, beaconState.SetCurrentSyncCommittee(sCom))
	tState := beaconState.Copy()
	blk, err := testutil.GenerateFullBlockAltair(tState, privKeys,
		&testutil.BlockGenConfig{NumAttestations: 1, NumVoluntaryExits: 0, NumDeposits: 0}, 1)
	require.NoError(t, err)

	return beaconState, blk
}
