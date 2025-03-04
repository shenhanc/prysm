package v2

import (
	"strconv"
	"sync"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestValidatorMap_DistinctCopy(t *testing.T) {
	count := uint64(100)
	vals := make([]*ethpb.Validator, 0, count)
	for i := uint64(1); i < count; i++ {
		someRoot := [32]byte{}
		someKey := [48]byte{}
		copy(someRoot[:], strconv.Itoa(int(i)))
		copy(someKey[:], strconv.Itoa(int(i)))
		vals = append(vals, &ethpb.Validator{
			PublicKey:                  someKey[:],
			WithdrawalCredentials:      someRoot[:],
			EffectiveBalance:           params.BeaconConfig().MaxEffectiveBalance,
			Slashed:                    false,
			ActivationEligibilityEpoch: 1,
			ActivationEpoch:            1,
			ExitEpoch:                  1,
			WithdrawableEpoch:          1,
		})
	}
	handler := stateutil.NewValMapHandler(vals)
	newHandler := handler.Copy()
	wantedPubkey := strconv.Itoa(22)
	handler.Set(bytesutil.ToBytes48([]byte(wantedPubkey)), 27)
	val1, _ := handler.Get(bytesutil.ToBytes48([]byte(wantedPubkey)))
	val2, _ := newHandler.Get(bytesutil.ToBytes48([]byte(wantedPubkey)))
	assert.NotEqual(t, val1, val2, "Values are supposed to be unequal due to copy")
}

func TestInitializeFromProto(t *testing.T) {
	type test struct {
		name  string
		state *statepb.BeaconStateAltair
		error string
	}
	initTests := []test{
		{
			name:  "nil state",
			state: nil,
			error: "received nil state",
		},
		{
			name: "nil validators",
			state: &statepb.BeaconStateAltair{
				Slot:       4,
				Validators: nil,
			},
		},
		{
			name:  "empty state",
			state: &statepb.BeaconStateAltair{},
		},
		// TODO: Add full state. Blocked by testutil migration.
	}
	for _, tt := range initTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := InitializeFromProto(tt.state)
			if tt.error != "" {
				require.ErrorContains(t, tt.error, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBeaconState_NoDeadlock(t *testing.T) {
	count := uint64(100)
	vals := make([]*ethpb.Validator, 0, count)
	for i := uint64(1); i < count; i++ {
		someRoot := [32]byte{}
		someKey := [48]byte{}
		copy(someRoot[:], strconv.Itoa(int(i)))
		copy(someKey[:], strconv.Itoa(int(i)))
		vals = append(vals, &ethpb.Validator{
			PublicKey:                  someKey[:],
			WithdrawalCredentials:      someRoot[:],
			EffectiveBalance:           params.BeaconConfig().MaxEffectiveBalance,
			Slashed:                    false,
			ActivationEligibilityEpoch: 1,
			ActivationEpoch:            1,
			ExitEpoch:                  1,
			WithdrawableEpoch:          1,
		})
	}
	st, err := InitializeFromProtoUnsafe(&statepb.BeaconStateAltair{
		Validators: vals,
	})
	assert.NoError(t, err)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		// Continuously lock and unlock the state
		// by acquiring the lock.
		for i := 0; i < 1000; i++ {
			for _, f := range st.stateFieldLeaves {
				f.Lock()
				if len(f.fieldLayers) == 0 {
					f.fieldLayers = make([][]*[32]byte, 10)
				}
				f.Unlock()
				f.reference.AddRef()
			}
		}
		wg.Done()
	}()
	// Constantly read from the offending portion
	// of the code to ensure there is no possible
	// recursive read locking.
	for i := 0; i < 1000; i++ {
		go func() {
			_ = st.FieldReferencesCount()
		}()
	}
	// Test will not terminate in the event of a deadlock.
	wg.Wait()
}

func TestInitializeFromProtoUnsafe(t *testing.T) {
	type test struct {
		name  string
		state *statepb.BeaconStateAltair
		error string
	}
	initTests := []test{
		{
			name:  "nil state",
			state: nil,
			error: "received nil state",
		},
		{
			name: "nil validators",
			state: &statepb.BeaconStateAltair{
				Slot:       4,
				Validators: nil,
			},
		},
		{
			name:  "empty state",
			state: &statepb.BeaconStateAltair{},
		},
		// TODO: Add full state. Blocked by testutil migration.
	}
	_ = initTests
}
