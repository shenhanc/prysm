package operations

import (
	"context"
	"path"
	"testing"

	"github.com/golang/snappy"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	prysmv2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
	"github.com/prysmaticlabs/prysm/proto/prysm/v2/block"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/spectest/utils"
)

func RunAttesterSlashingTest(t *testing.T, config string) {
	require.NoError(t, utils.SetConfig(t, config))
	testFolders, testsFolderPath := utils.TestFolders(t, config, "altair", "operations/attester_slashing/pyspec_tests")
	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			folderPath := path.Join(testsFolderPath, folder.Name())
			attSlashingFile, err := testutil.BazelFileBytes(folderPath, "attester_slashing.ssz_snappy")
			require.NoError(t, err)
			attSlashingSSZ, err := snappy.Decode(nil /* dst */, attSlashingFile)
			require.NoError(t, err, "Failed to decompress")
			attSlashing := &ethpb.AttesterSlashing{}
			require.NoError(t, attSlashing.UnmarshalSSZ(attSlashingSSZ), "Failed to unmarshal")

			body := &prysmv2.BeaconBlockBodyAltair{AttesterSlashings: []*ethpb.AttesterSlashing{attSlashing}}
			RunBlockOperationTest(t, folderPath, body, func(ctx context.Context, s state.BeaconState, b block.SignedBeaconBlock) (state.BeaconState, error) {
				return blocks.ProcessAttesterSlashings(ctx, s, b.Block().Body().AttesterSlashings(), altair.SlashValidator)
			})
		})
	}
}