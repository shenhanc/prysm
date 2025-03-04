package testutil

import (
	"testing"

	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestHydrateAttestation(t *testing.T) {
	a := HydrateAttestation(&ethpb.Attestation{})
	_, err := a.HashTreeRoot()
	require.NoError(t, err)
	require.DeepEqual(t, a.Signature, make([]byte, params.BeaconConfig().BLSSignatureLength))
}

func TestHydrateAttestationData(t *testing.T) {
	d := HydrateAttestationData(&ethpb.AttestationData{})
	_, err := d.HashTreeRoot()
	require.NoError(t, err)
	require.DeepEqual(t, d.BeaconBlockRoot, make([]byte, 32))
	require.DeepEqual(t, d.Target.Root, make([]byte, 32))
	require.DeepEqual(t, d.Source.Root, make([]byte, 32))
}

func TestHydrateIndexedAttestation(t *testing.T) {
	a := &ethpb.IndexedAttestation{}
	a = HydrateIndexedAttestation(a)
	_, err := a.HashTreeRoot()
	require.NoError(t, err)
	_, err = a.Data.HashTreeRoot()
	require.NoError(t, err)
}

func TestGenerateAttestations_EpochBoundary(t *testing.T) {
	gs, pk := DeterministicGenesisState(t, 32)
	_, err := GenerateAttestations(gs, pk, 1, params.BeaconConfig().SlotsPerEpoch, false)
	require.NoError(t, err)
}
