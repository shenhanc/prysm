package genesis

import (
	_ "embed"

	"github.com/golang/snappy"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/params"
)

var (
	//go:embed mainnet.ssz.snappy
	mainnetRawSSZCompressed []byte // 1.8Mb
)

// State returns a copy of the genesis state from a hardcoded value.
func State(name string) (*state.BeaconState, error) {
	switch name {
	case params.ConfigNames[params.Mainnet]:
		return load(mainnetRawSSZCompressed)
	default:
		// No state found.
		return nil, nil
	}
}

// load a compressed ssz state file into a beacon state struct.
func load(b []byte) (*state.BeaconState, error) {
	st := &statepb.BeaconState{}
	b, err := snappy.Decode(nil /*dst*/, b)
	if err != nil {
		return nil, err
	}
	if err := st.UnmarshalSSZ(b); err != nil {
		return nil, err
	}
	return state.InitializeFromProtoUnsafe(st)
}
