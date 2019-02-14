package http2

import (
	"encoding/binary"
	"io"
)

// FrameHeaderLen frame header length
const FrameHeaderLen = 9

// FrameType frame type
type FrameType uint8

const (
	// FrameData body splice
	FrameData FrameType = 0x0
	// FrameHeaders header start
	FrameHeaders FrameType = 0x1
	// FramePriority priority
	FramePriority FrameType = 0x2
	// FrameRSTStream rst stream
	FrameRSTStream FrameType = 0x3
	// FrameSettings settings frame
	FrameSettings FrameType = 0x4
	// FramePushPromise push promise
	FramePushPromise FrameType = 0x5
	// FramePing ping
	FramePing FrameType = 0x6
	// FrameGoAway goaway
	FrameGoAway FrameType = 0x7
	// FrameWindowUpdate data flow
	FrameWindowUpdate FrameType = 0x8
	// FrameContinuation continuation
	FrameContinuation FrameType = 0x9
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

// ReadFrame 从流中读取一个 frame
func ReadFrame(r io.Reader) (*FrameHeader, []byte, error) {
	var buf = frameCache.Get().(*[]byte)
	header, err := readFrameHeader(*buf, r)
	frameCache.Put(buf)
	if err != nil {
		return nil, nil, err
	}
	payload := make([]byte, header.Length)
	_, err = io.ReadFull(r, payload)
	return header, payload, err
}
