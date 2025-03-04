package metadata

import (
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/go-bitfield"
	pb "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)

// Metadata returns the interface of a p2p metadata type.
type Metadata interface {
	SequenceNumber() uint64
	AttnetsBitfield() bitfield.Bitvector64
	InnerObject() interface{}
	IsNil() bool
	Copy() Metadata
	ssz.Marshaler
	ssz.Unmarshaler
	MetadataObjV0() *pb.MetaDataV0
	MetadataObjV1() *pb.MetaDataV1
	Version() int
}
