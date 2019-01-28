package http2

import (
	"encoding/binary"
	"io"
)

// FrameType frame type
type FrameType uint8

const (
	FrameData         FrameType = 0x0
	FrameHeaders      FrameType = 0x1
	FramePriority     FrameType = 0x2
	FrameRSTStream    FrameType = 0x3
	FrameSettings     FrameType = 0x4
	FramePushPromise  FrameType = 0x5
	FramePing         FrameType = 0x6
	FrameGoAway       FrameType = 0x7
	FrameWindowUpdate FrameType = 0x8
	FrameContinuation FrameType = 0x9
)

// Flags flag type
type Flags uint8

func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

const (
	// Data Frame

	FlagDataEndStream Flags = 0x1
	FlagDataPadded    Flags = 0x8

	// Headers Frame

	FlagHeadersEndStream  Flags = 0x1
	FlagHeadersEndHeaders Flags = 0x4
	FlagHeadersPadded     Flags = 0x8
	FlagHeadersPriority   Flags = 0x20

	// Settings Frame

	FlagSettingsAck Flags = 0x1

	// Ping Frame

	FlagPingAck Flags = 0x1

	// Continuation Frame

	FlagContinuationEndHeaders Flags = 0x4

	FlagPushPromiseEndHeaders Flags = 0x4
	FlagPushPromisePadded     Flags = 0x8
)

// PadBit (1<<31) - 1
const PadBit uint32 = 0x7fffffff

// FrameHeader frame header struct
type FrameHeader struct {
	Type     FrameType
	Flags    Flags
	Length   uint32
	StreamID uint32
}

// ReadFrameHeader 从流中解析出 frame header
func readFrameHeader(buf []byte, conn io.Reader) (frame *FrameHeader, err error) {
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	frame = &FrameHeader{
		// 把长度解析为 uint32, 由于长度不足 32 bit 手动通过位操作
		Length: (uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])),
		Type:   FrameType(buf[3]),
		// flags 也为 uint8
		Flags: Flags(buf[4]),
		// StreamID 正好长度为 32bit 通过binary.BigEndian.Uint32进行转换。
		StreamID: binary.BigEndian.Uint32(buf[5:]) & PadBit,
	}
	return frame, nil
}
