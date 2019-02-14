package http2

import (
	"encoding/binary"
	"io"
	"sync"
)

func readByte(p []byte) ([]byte, byte, error) {
	if len(p) == 0 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	return p[1:], p[0], nil
}

func readUint32(p []byte) ([]byte, uint32, error) {
	if len(p) < 4 {
		return nil, 0, io.ErrUnexpectedEOF
	}
	return p[4:], binary.BigEndian.Uint32(p[:4]), nil
}

var frameCache = &sync.Pool{
	New: func() interface{} {
		buf := make([]byte, FrameHeaderLen)
		return &buf
	},
}

