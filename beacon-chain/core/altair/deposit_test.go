package altair_test

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	stateAltair "github.com/prysmaticlabs/prysm/beacon-chain/state/v2"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/shared/trieutil"
)

func TestProcessDeposits_SameValidatorMultipleDepositsSameBlock(t *testing.T) {
	// Same validator created 3 valid deposits within the same block
	dep, _, err := testutil.DeterministicDepositsAndKeysSameValidator(3)
	require.NoError(t, err)
	eth1Data, err := testutil.DeterministicEth1Data(len(dep))
	require.NoError(t, err)
	registry := []*ethpb.Validator{
		{
			PublicKey:             []byte{1},
			WithdrawalCredentials: []byte{1, 2, 3},
		},
	}
	balances := []uint64{0}
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Validators: registry,
		Balances:   balances,
		Eth1Data:   eth1Data,
		Fork: &statepb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
		},
	})
	require.NoError(t, err)
	newState, err := altair.ProcessDeposits(context.Background(), beaconState, []*ethpb.Deposit{dep[0], dep[1], dep[2]})
	require.NoError(t, err, "Expected block deposits to process correctly")
	require.Equal(t, 2, len(newState.Validators()), "Incorrect validator count")
}

func TestProcessDeposits_MerkleBranchFailsVerification(t *testing.T) {
	deposit := &ethpb.Deposit{
		Data: &ethpb.Deposit_Data{
			PublicKey:             bytesutil.PadTo([]byte{1, 2, 3}, 48),
			WithdrawalCredentials: make([]byte, 32),
			Signature:             make([]byte, 96),
		},
	}
	leaf, err := deposit.Data.HashTreeRoot()
	require.NoError(t, err)

	// We then create a merkle branch for the test.
	depositTrie, err := trieutil.GenerateTrieFromItems([][]byte{leaf[:]}, params.BeaconConfig().DepositContractTreeDepth)
	require.NoError(t, err, "Could not generate trie")
	proof, err := depositTrie.MerkleProof(0)
	require.NoError(t, err, "Could not generate proof")

	deposit.Proof = proof
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Eth1Data: &ethpb.Eth1Data{
			DepositRoot: []byte{0},
			BlockHash:   []byte{1},
		},
	})
	require.NoError(t, err)
	want := "deposit root did not verify"
	_, err = altair.ProcessDeposits(context.Background(), beaconState, []*ethpb.Deposit{deposit})
	require.ErrorContains(t, want, err)
}

func TestProcessDeposits_AddsNewValidatorDeposit(t *testing.T) {
	dep, _, err := testutil.DeterministicDepositsAndKeys(1)
	require.NoError(t, err)
	eth1Data, err := testutil.DeterministicEth1Data(len(dep))
	require.NoError(t, err)

	registry := []*ethpb.Validator{
		{
			PublicKey:             []byte{1},
			WithdrawalCredentials: []byte{1, 2, 3},
		},
	}
	balances := []uint64{0}
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Validators: registry,
		Balances:   balances,
		Eth1Data:   eth1Data,
		Fork: &statepb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
		},
	})
	require.NoError(t, err)
	newState, err := altair.ProcessDeposits(context.Background(), beaconState, []*ethpb.Deposit{dep[0]})
	require.NoError(t, err, "Expected block deposits to process correctly")
	if newState.Balances()[1] != dep[0].Data.Amount {
		t.Errorf(
			"Expected state validator balances index 0 to equal %d, received %d",
			dep[0].Data.Amount,
			newState.Balances()[1],
		)
	}
}

func TestProcessDeposits_RepeatedDeposit_IncreasesValidatorBalance(t *testing.T) {
	sk, err := bls.RandKey()
	require.NoError(t, err)
	deposit := &ethpb.Deposit{
		Data: &ethpb.Deposit_Data{
			PublicKey:             sk.PublicKey().Marshal(),
			Amount:                1000,
			WithdrawalCredentials: make([]byte, 32),
			Signature:             make([]byte, 96),
		},
	}
	sr, err := helpers.ComputeSigningRoot(deposit.Data, bytesutil.ToBytes(3, 32))
	require.NoError(t, err)
	sig := sk.Sign(sr[:])
	deposit.Data.Signature = sig.Marshal()
	leaf, err := deposit.Data.HashTreeRoot()
	require.NoError(t, err)

	// We then create a merkle branch for the test.
	depositTrie, err := trieutil.GenerateTrieFromItems([][]byte{leaf[:]}, params.BeaconConfig().DepositContractTreeDepth)
	require.NoError(t, err, "Could not generate trie")
	proof, err := depositTrie.MerkleProof(0)
	require.NoError(t, err, "Could not generate proof")

	deposit.Proof = proof
	registry := []*ethpb.Validator{
		{
			PublicKey: []byte{1, 2, 3},
		},
		{
			PublicKey:             sk.PublicKey().Marshal(),
			WithdrawalCredentials: []byte{1},
		},
	}
	balances := []uint64{0, 50}
	root := depositTrie.Root()
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Validators: registry,
		Balances:   balances,
		Eth1Data: &ethpb.Eth1Data{
			DepositRoot: root[:],
			BlockHash:   root[:],
		},
	})
	require.NoError(t, err)
	newState, err := altair.ProcessDeposits(context.Background(), beaconState, []*ethpb.Deposit{deposit})
	require.NoError(t, err, "Process deposit failed")
	require.Equal(t, uint64(1000+50), newState.Balances()[1], "Expected balance at index 1 to be 1050")
}

func TestProcessDeposit_AddsNewValidatorDeposit(t *testing.T) {
	// Similar to TestProcessDeposits_AddsNewValidatorDeposit except that this test directly calls ProcessDeposit
	dep, _, err := testutil.DeterministicDepositsAndKeys(1)
	require.NoError(t, err)
	eth1Data, err := testutil.DeterministicEth1Data(len(dep))
	require.NoError(t, err)

	registry := []*ethpb.Validator{
		{
			PublicKey:             []byte{1},
			WithdrawalCredentials: []byte{1, 2, 3},
		},
	}
	balances := []uint64{0}
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Validators: registry,
		Balances:   balances,
		Eth1Data:   eth1Data,
		Fork: &statepb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
		},
	})
	require.NoError(t, err)
	newState, err := altair.ProcessDeposit(context.Background(), beaconState, dep[0], true)
	require.NoError(t, err, "Process deposit failed")
	require.Equal(t, 2, len(newState.Validators()), "Expected validator list to have length 2")
	require.Equal(t, 2, len(newState.Balances()), "Expected validator balances list to have length 2")
	if newState.Balances()[1] != dep[0].Data.Amount {
		t.Errorf(
			"Expected state validator balances index 1 to equal %d, received %d",
			dep[0].Data.Amount,
			newState.Balances()[1],
		)
	}
}

func TestProcessDeposit_SkipsInvalidDeposit(t *testing.T) {
	// Same test settings as in TestProcessDeposit_AddsNewValidatorDeposit, except that we use an invalid signature
	dep, _, err := testutil.DeterministicDepositsAndKeys(1)
	require.NoError(t, err)
	dep[0].Data.Signature = make([]byte, 96)
	trie, _, err := testutil.DepositTrieFromDeposits(dep)
	require.NoError(t, err)
	root := trie.Root()
	eth1Data := &ethpb.Eth1Data{
		DepositRoot:  root[:],
		DepositCount: 1,
	}
	registry := []*ethpb.Validator{
		{
			PublicKey:             []byte{1},
			WithdrawalCredentials: []byte{1, 2, 3},
		},
	}
	balances := []uint64{0}
	beaconState, err := stateAltair.InitializeFromProto(&statepb.BeaconStateAltair{
		Validators: registry,
		Balances:   balances,
		Eth1Data:   eth1Data,
		Fork: &statepb.Fork{
			PreviousVersion: params.BeaconConfig().GenesisForkVersion,
			CurrentVersion:  params.BeaconConfig().GenesisForkVersion,
		},
	})
	require.NoError(t, err)
	newState, err := altair.ProcessDeposit(context.Background(), beaconState, dep[0], true)
	require.NoError(t, err, "Expected invalid block deposit to be ignored without error")

	if newState.Eth1DepositIndex() != 1 {
		t.Errorf(
			"Expected Eth1DepositIndex to be increased by 1 after processing an invalid deposit, received change: %v",
			newState.Eth1DepositIndex(),
		)
	}
	if len(newState.Validators()) != 1 {
		t.Errorf("Expected validator list to have length 1, received: %v", len(newState.Validators()))
	}
	if len(newState.Balances()) != 1 {
		t.Errorf("Expected validator balances list to have length 1, received: %v", len(newState.Balances()))
	}
	if newState.Balances()[0] != 0 {
		t.Errorf("Expected validator balance at index 0 to stay 0, received: %v", newState.Balances()[0])
	}
}
