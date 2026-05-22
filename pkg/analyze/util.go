package analyze

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
)

func readDelimited(r io.Reader, m proto.Message) (n int, pos uint64, err error) {
	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	var headerBuf [binary.MaxVarintLen32]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 { // i.e. no varint has been decoded yet.
		if bytesRead >= len(headerBuf) {
			return bytesRead, pos, errors.New("invalid varint32 encountered")
		}
		// We have to read byte by byte here to avoid reading more bytes
		// than required. Each read byte is appended to what we have
		// read before.
		newBytesRead, err := r.Read(headerBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if err != nil {
				return bytesRead, pos, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		// Now present everything read so far to the varint decoder and
		// see if a varint can be decoded already.
		messageLength, varIntBytes = proto.DecodeVarint(headerBuf[:bytesRead])
	}

	messageBuf := make([]byte, messageLength)
	newBytesRead, err := io.ReadFull(r, messageBuf)
	bytesRead += newBytesRead
	if err != nil {
		return bytesRead, pos, err
	}
	pos = messageLength + uint64(varIntBytes)
	return bytesRead, pos, proto.Unmarshal(messageBuf, m)
}
