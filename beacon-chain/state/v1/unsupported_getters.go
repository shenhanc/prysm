package v1

import (
	"github.com/pkg/errors"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
)

// CurrentEpochParticipation is not supported for phase 0 beacon state.
func (b *BeaconState) CurrentEpochParticipation() ([]byte, error) {
	return nil, errors.New("CurrentEpochParticipation is not supported for phase 0 beacon state")
}

// PreviousEpochParticipation is not supported for phase 0 beacon state.
func (b *BeaconState) PreviousEpochParticipation() ([]byte, error) {
	return nil, errors.New("PreviousEpochParticipation is not supported for phase 0 beacon state")
}

// InactivityScores is not supported for phase 0 beacon state.
func (b *BeaconState) InactivityScores() ([]uint64, error) {
	return nil, errors.New("InactivityScores is not supported for phase 0 beacon state")
}

// CurrentSyncCommittee is not supported for phase 0 beacon state.
func (b *BeaconState) CurrentSyncCommittee() (*statepb.SyncCommittee, error) {
	return nil, errors.New("CurrentSyncCommittee is not supported for phase 0 beacon state")
}

// NextSyncCommittee is not supported for phase 0 beacon state.
func (b *BeaconState) NextSyncCommittee() (*statepb.SyncCommittee, error) {
	return nil, errors.New("NextSyncCommittee is not supported for phase 0 beacon state")
}
