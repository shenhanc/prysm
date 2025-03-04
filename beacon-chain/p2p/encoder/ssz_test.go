package encoder_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"testing"

	gogo "github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p/encoder"
	statepb "github.com/prysmaticlabs/prysm/proto/prysm/v2/state"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"google.golang.org/protobuf/proto"
)

func TestSszNetworkEncoder_RoundTrip(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	testRoundTripWithLength(t, e)
	testRoundTripWithGossip(t, e)
}

func TestSszNetworkEncoder_FailsSnappyLength(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	att := &statepb.Fork{}
	data := make([]byte, 32)
	binary.PutUvarint(data, encoder.MaxGossipSize+32)
	err := e.DecodeGossip(data, att)
	require.ErrorContains(t, "snappy message exceeds max size", err)
}

func testRoundTripWithLength(t *testing.T, e *encoder.SszNetworkEncoder) {
	buf := new(bytes.Buffer)
	msg := &statepb.Fork{
		PreviousVersion: []byte("fooo"),
		CurrentVersion:  []byte("barr"),
		Epoch:           9001,
	}
	_, err := e.EncodeWithMaxLength(buf, msg)
	require.NoError(t, err)
	decoded := &statepb.Fork{}
	require.NoError(t, e.DecodeWithMaxLength(buf, decoded))
	if !proto.Equal(decoded, msg) {
		t.Logf("decoded=%+v\n", decoded)
		t.Error("Decoded message is not the same as original")
	}
}

func testRoundTripWithGossip(t *testing.T, e *encoder.SszNetworkEncoder) {
	buf := new(bytes.Buffer)
	msg := &statepb.Fork{
		PreviousVersion: []byte("fooo"),
		CurrentVersion:  []byte("barr"),
		Epoch:           9001,
	}
	_, err := e.EncodeGossip(buf, msg)
	require.NoError(t, err)
	decoded := &statepb.Fork{}
	require.NoError(t, e.DecodeGossip(buf.Bytes(), decoded))
	if !proto.Equal(decoded, msg) {
		t.Logf("decoded=%+v\n", decoded)
		t.Error("Decoded message is not the same as original")
	}
}

func TestSszNetworkEncoder_EncodeWithMaxLength(t *testing.T) {
	buf := new(bytes.Buffer)
	msg := &statepb.Fork{
		PreviousVersion: []byte("fooo"),
		CurrentVersion:  []byte("barr"),
		Epoch:           9001,
	}
	e := &encoder.SszNetworkEncoder{}
	params.SetupTestConfigCleanup(t)
	c := params.BeaconNetworkConfig()
	c.MaxChunkSize = uint64(5)
	params.OverrideBeaconNetworkConfig(c)
	_, err := e.EncodeWithMaxLength(buf, msg)
	wanted := fmt.Sprintf("which is larger than the provided max limit of %d", params.BeaconNetworkConfig().MaxChunkSize)
	assert.ErrorContains(t, wanted, err)
}

func TestSszNetworkEncoder_DecodeWithMaxLength(t *testing.T) {
	buf := new(bytes.Buffer)
	msg := &statepb.Fork{
		PreviousVersion: []byte("fooo"),
		CurrentVersion:  []byte("barr"),
		Epoch:           4242,
	}
	e := &encoder.SszNetworkEncoder{}
	params.SetupTestConfigCleanup(t)
	c := params.BeaconNetworkConfig()
	maxChunkSize := uint64(5)
	c.MaxChunkSize = maxChunkSize
	params.OverrideBeaconNetworkConfig(c)
	_, err := e.EncodeGossip(buf, msg)
	require.NoError(t, err)
	decoded := &statepb.Fork{}
	err = e.DecodeWithMaxLength(buf, decoded)
	wanted := fmt.Sprintf("goes over the provided max limit of %d", maxChunkSize)
	assert.ErrorContains(t, wanted, err)
}

func TestSszNetworkEncoder_DecodeWithMultipleFrames(t *testing.T) {
	buf := new(bytes.Buffer)
	st, _ := testutil.DeterministicGenesisState(t, 100)
	e := &encoder.SszNetworkEncoder{}
	params.SetupTestConfigCleanup(t)
	c := params.BeaconNetworkConfig()
	// 4 * 1 Mib
	maxChunkSize := uint64(1 << 22)
	c.MaxChunkSize = maxChunkSize
	params.OverrideBeaconNetworkConfig(c)
	_, err := e.EncodeWithMaxLength(buf, st.InnerStateUnsafe().(*statepb.BeaconState))
	require.NoError(t, err)
	// Max snappy block size
	if buf.Len() <= 76490 {
		t.Errorf("buffer smaller than expected, wanted > %d but got %d", 76490, buf.Len())
	}
	decoded := new(statepb.BeaconState)
	err = e.DecodeWithMaxLength(buf, decoded)
	assert.NoError(t, err)
}
func TestSszNetworkEncoder_NegativeMaxLength(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	length, err := e.MaxLength(0xfffffffffff)

	assert.Equal(t, 0, length, "Received non zero length on bad message length")
	assert.ErrorContains(t, "max encoded length is negative", err)
}

func TestSszNetworkEncoder_MaxInt64(t *testing.T) {
	e := &encoder.SszNetworkEncoder{}
	length, err := e.MaxLength(math.MaxInt64 + 1)

	assert.Equal(t, 0, length, "Received non zero length on bad message length")
	assert.ErrorContains(t, "invalid length provided", err)
}

func TestSszNetworkEncoder_DecodeWithBadSnappyStream(t *testing.T) {
	st := newBadSnappyStream()
	e := &encoder.SszNetworkEncoder{}
	decoded := new(statepb.Fork)
	err := e.DecodeWithMaxLength(st, decoded)
	assert.ErrorContains(t, io.EOF.Error(), err)
}

type badSnappyStream struct {
	varint []byte
	header []byte
	repeat []byte
	i      int
	// count how many times it was read
	counter int
	// count bytes read so far
	total int
}

func newBadSnappyStream() *badSnappyStream {
	const (
		magicBody  = "sNaPpY"
		magicChunk = "\xff\x06\x00\x00" + magicBody
	)

	header := make([]byte, len(magicChunk))
	// magicChunk == chunkTypeStreamIdentifier byte ++ 3 byte little endian len(magic body) ++ 6 byte magic body

	// header is a special chunk type, with small fixed length, to add some magic to claim it's really snappy.
	copy(header, magicChunk) // snappy library constants help us construct the common header chunk easily.

	payload := make([]byte, 4)

	// byte 0 is chunk type
	// Exploit any fancy ignored chunk type
	//   Section 4.4 Padding (chunk type 0xfe).
	//   Section 4.6. Reserved skippable chunks (chunk types 0x80-0xfd).
	payload[0] = 0xfe

	// byte 1,2,3 are chunk length (little endian)
	payload[1] = 0
	payload[2] = 0
	payload[3] = 0

	return &badSnappyStream{
		varint:  gogo.EncodeVarint(1000),
		header:  header,
		repeat:  payload,
		i:       0,
		counter: 0,
		total:   0,
	}
}

func (b *badSnappyStream) Read(p []byte) (n int, err error) {
	// Stream out varint bytes first to make test happy.
	if len(b.varint) > 0 {
		copy(p, b.varint[:1])
		b.varint = b.varint[1:]
		return 1, nil
	}
	defer func() {
		b.counter += 1
		b.total += n
	}()
	if len(b.repeat) == 0 {
		panic("no bytes to repeat")
	}
	if len(b.header) > 0 {
		n = copy(p, b.header)
		b.header = b.header[n:]
		return
	}
	for n < len(p) {
		n += copy(p[n:], b.repeat[b.i:])
		b.i = (b.i + n) % len(b.repeat)
	}
	return
}
