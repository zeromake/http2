package http2

import "encoding/binary"

// RSTStreamFrame rst stream frame
type RSTStreamFrame struct {
	FrameHeader
	ErrCode ErrCode
}

// ParseRSTStreamFrame parse rst stream frame
func ParseRSTStreamFrame(fh *FrameHeader, p []byte) (*RSTStreamFrame, error) {
	if len(p) != 4 {
		return nil, ConnectionError(ErrCodeFrameSize)
	}
	if fh.StreamID == 0 {
		return nil, ConnectionError(ErrCodeProtocol)
	}
	return &RSTStreamFrame{*fh, ErrCode(binary.BigEndian.Uint32(p[:4]))}, nil
}

