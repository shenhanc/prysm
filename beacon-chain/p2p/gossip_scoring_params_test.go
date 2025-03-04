package p2p

import (
	"context"
	"testing"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	dbutil "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestCorrect_ActiveValidatorsCount(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	cfg := params.MainnetConfig()
	cfg.ConfigName = "test"

	params.OverrideBeaconConfig(cfg)

	db := dbutil.SetupDB(t)
	s := &Service{
		ctx: context.Background(),
		cfg: &Config{DB: db},
	}
	bState, err := testutil.NewBeaconState(func(state *statepb.BeaconState) error {
		validators := make([]*ethpb.Validator, params.BeaconConfig().MinGenesisActiveValidatorCount)
		for i := 0; i < len(validators); i++ {
			validators[i] = &ethpb.Validator{
				PublicKey:             make([]byte, 48),
				WithdrawalCredentials: make([]byte, 32),
				ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
				Slashed:               false,
			}
		}
		state.Validators = validators
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, db.SaveGenesisData(s.ctx, bState))

	vals, err := s.retrieveActiveValidators()
	assert.NoError(t, err, "genesis state not retrieved")
	assert.Equal(t, int(params.BeaconConfig().MinGenesisActiveValidatorCount), int(vals), "mainnet genesis active count isn't accurate")
	for i := 0; i < 100; i++ {
		require.NoError(t, bState.AppendValidator(&ethpb.Validator{
			PublicKey:             make([]byte, 48),
			WithdrawalCredentials: make([]byte, 32),
			ExitEpoch:             params.BeaconConfig().FarFutureEpoch,
			Slashed:               false,
		}))
	}
	require.NoError(t, bState.SetSlot(10000))
	require.NoError(t, db.SaveState(s.ctx, bState, [32]byte{'a'}))
	// Reset count
	s.activeValidatorCount = 0

	// Retrieve last archived state.
	vals, err = s.retrieveActiveValidators()
	assert.NoError(t, err, "genesis state not retrieved")
	assert.Equal(t, int(params.BeaconConfig().MinGenesisActiveValidatorCount)+100, int(vals), "mainnet genesis active count isn't accurate")
}

func TestLoggingParameters(t *testing.T) {
	logGossipParameters("testing", nil)
	logGossipParameters("testing", &pubsub.TopicScoreParams{})
	// Test out actual gossip parameters.
	logGossipParameters("testing", defaultBlockTopicParams())
	p, err := defaultAggregateSubnetTopicParams(10000)
	assert.NoError(t, err)
	logGossipParameters("testing", p)
	p, err = defaultAggregateTopicParams(10000)
	assert.NoError(t, err)
	logGossipParameters("testing", p)
	logGossipParameters("testing", defaultAttesterSlashingTopicParams())
	logGossipParameters("testing", defaultProposerSlashingTopicParams())
	logGossipParameters("testing", defaultVoluntaryExitTopicParams())
}
