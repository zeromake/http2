package http2

import (
	"encoding/binary"
	"fmt"
)

// PriorityParam priority param
type PriorityParam struct {
	StreamDep uint32
	Exclusive bool
	Weight    uint8
}

// PriorityFrame priority frame
type PriorityFrame struct {
	FrameHeader
	PriorityParam
}

// ParsePriorityFrame parse priority frame
func ParsePriorityFrame(fh *FrameHeader, payload []byte) (*PriorityFrame, error) {
	if fh.StreamID == 0 {
		return nil, connError{
			ErrCodeProtocol,
			"PRIORITY frame with stream ID 0",
		}
	}
	var payloadLength = len(payload)
	if payloadLength != 5 {
		return nil, connError{
			ErrCodeFrameSize,
			fmt.Sprintf(
				"PRIORITY frame payload size was %d; want 5",
				payloadLength,
			),
		}
	}
	v := binary.BigEndian.Uint32(payload[:4])
	// E 不处理
	streamID := v & PadBit // mask off high bit
	return &PriorityFrame{
		FrameHeader: *fh,
		PriorityParam: PriorityParam{
			Weight:    payload[4],
			StreamDep: streamID,
			Exclusive: streamID != v, // was high bit set?
		},
	}, nil
}
