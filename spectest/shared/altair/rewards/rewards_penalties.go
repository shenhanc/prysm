package rewards

import (
	"context"
	"encoding/binary"
	"fmt"
	"path"
	"reflect"
	"testing"

	"github.com/golang/snappy"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	stateAltair "github.com/prysmaticlabs/prysm/beacon-chain/state/v2"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/spectest/utils"
)

// Delta contains list of rewards and penalties.
type Delta struct {
	Rewards   []uint64 `json:"rewards"`
	Penalties []uint64 `json:"penalties"`
}

// unmarshalSSZ deserializes specs data into a simple aggregating container.
func (d *Delta) unmarshalSSZ(buf []byte) error {
	offset1 := binary.LittleEndian.Uint32(buf[:4])
	offset2 := binary.LittleEndian.Uint32(buf[4:8])

	for i := uint32(0); i < offset2-offset1; i += 8 {
		d.Rewards = append(d.Rewards, binary.LittleEndian.Uint64(buf[offset1+i:offset1+i+8]))
		d.Penalties = append(d.Penalties, binary.LittleEndian.Uint64(buf[offset2+i:offset2+i+8]))
	}
	return nil
}

var deltaFiles = []string{
	"source_deltas.ssz_snappy",
	"target_deltas.ssz_snappy",
	"head_deltas.ssz_snappy",
	"inactivity_penalty_deltas.ssz_snappy",
}

// RunPrecomputeRewardsAndPenaltiesTests executes "rewards/{basic, leak, random}" tests.
func RunPrecomputeRewardsAndPenaltiesTests(t *testing.T, config string) {
	require.NoError(t, utils.SetConfig(t, config))
	testTypes := []string{"basic", "leak", "random"}
	for _, testType := range testTypes {
		testFolders, testsFolderPath := utils.TestFolders(t, config, "altair", fmt.Sprintf("rewards/%s/pyspec_tests", testType))
		for _, folder := range testFolders {
			helpers.ClearCache()
			t.Run(fmt.Sprintf("%v/%v", testType, folder.Name()), func(t *testing.T) {
				folderPath := path.Join(testsFolderPath, folder.Name())
				runPrecomputeRewardsAndPenaltiesTest(t, folderPath)
			})
		}
	}
}

func runPrecomputeRewardsAndPenaltiesTest(t *testing.T, testFolderPath string) {
	ctx := context.Background()
	preBeaconStateFile, err := testutil.BazelFileBytes(path.Join(testFolderPath, "pre.ssz_snappy"))
	require.NoError(t, err)
	preBeaconStateSSZ, err := snappy.Decode(nil /* dst */, preBeaconStateFile)
	require.NoError(t, err, "Failed to decompress")
	preBeaconStateBase := &statepb.BeaconStateAltair{}
	require.NoError(t, preBeaconStateBase.UnmarshalSSZ(preBeaconStateSSZ), "Failed to unmarshal")
	preBeaconState, err := stateAltair.InitializeFromProto(preBeaconStateBase)
	require.NoError(t, err)

	vp, bp, err := altair.InitializeEpochValidators(ctx, preBeaconState)
	require.NoError(t, err)
	vp, bp, err = altair.ProcessEpochParticipation(ctx, preBeaconState, bp, vp)
	require.NoError(t, err)
	rewards, penalties, err := altair.AttestationsDelta(preBeaconState, bp, vp)
	require.NoError(t, err)

	totalSpecTestRewards := make([]uint64, len(rewards))
	totalSpecTestPenalties := make([]uint64, len(penalties))

	for _, dFile := range deltaFiles {
		sourceFile, err := testutil.BazelFileBytes(path.Join(testFolderPath, dFile))
		require.NoError(t, err)
		sourceSSZ, err := snappy.Decode(nil /* dst */, sourceFile)
		require.NoError(t, err, "Failed to decompress")
		d := &Delta{}
		require.NoError(t, d.unmarshalSSZ(sourceSSZ), "Failed to unmarshal")
		for i, reward := range d.Rewards {
			totalSpecTestRewards[i] += reward
		}
		for i, penalty := range d.Penalties {
			totalSpecTestPenalties[i] += penalty
		}
	}

	if !reflect.DeepEqual(rewards, totalSpecTestRewards) {
		t.Error("Rewards don't match")
		t.Log(rewards)
		t.Log(totalSpecTestRewards)
	}
	if !reflect.DeepEqual(penalties, totalSpecTestPenalties) {
		t.Error("Penalties don't match")
		t.Log(penalties)
		t.Log(totalSpecTestPenalties)
	}
}
