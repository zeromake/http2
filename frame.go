package http2

import (
	"encoding/binary"
	"fmt"
	"io"
)

const FrameHeaderLen = 9

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

// ReadFrame 从流中读取一个 frame
func ReadFrame(r io.Reader) (*FrameHeader, []byte, error) {
	// buf可以考虑通过池子复用
	var buf = frameCache.Get().([]byte)
	header, err := readFrameHeader(buf, r)
	frameCache.Put(buf)
	if err != nil {
		return nil, nil, err
	}
	payload := make([]byte, header.Length)
	_, err = io.ReadFull(r, payload)
	return header, payload, err
}

// DataFrame body
type DataFrame struct {
	FrameHeader
	Data []byte
}

// ParseDataFrame 解析
func ParseDataFrame(fh *FrameHeader, payload []byte) (*DataFrame, error) {
	if fh.StreamID == 0 {
		return nil, connError{ErrCodeProtocol, "DATA frame with stream ID 0"}
	}
	f := &DataFrame{FrameHeader: *fh}

	var padSize byte
	if fh.Flags.Has(FlagDataPadded) {
		var err error
		payload, padSize, err = readByte(payload)
		if err != nil {
			return nil, err
		}
	}
	if int(padSize) > len(payload) {
		return nil, connError{ErrCodeProtocol, "pad size larger than data payload"}
	}
	f.Data = payload[:len(payload)-int(padSize)]
	return f, nil
}

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

// SettingID setting id
type SettingID uint16

// Setting setting
type Setting struct {
	ID  SettingID
	Val uint32
}

// SettingFrame setting frame
type SettingFrame struct {
	FrameHeader
	Settings map[SettingID]Setting
}

// ParserSettings parse setting frame
func ParserSettings(header *FrameHeader, payload []byte) (*SettingFrame, error) {
	settings := map[SettingID]Setting{}
	num := len(payload) / 6
	for i := 0; i < num; i++ {
		id := SettingID(binary.BigEndian.Uint16(payload[i*6 : i*6+2]))
		s := Setting{
			ID:  id,
			Val: binary.BigEndian.Uint32(payload[i*6+2 : i*6+6]),
		}
		settings[id] = s
	}
	return &SettingFrame{
		*header,
		settings,
	}, nil
}

// PushPromiseFrame push promise frame
type PushPromiseFrame struct {
	FrameHeader
	PromiseID     uint32
	HeaderFragBuf []byte // not owned
}

// ParsePushPromise parse push promise
func ParsePushPromise(fh *FrameHeader, p []byte) (*PushPromiseFrame, error) {
	var err error
	pp := &PushPromiseFrame{
		FrameHeader: *fh,
	}
	if pp.StreamID == 0 {
		return nil, ConnectionError(ErrCodeProtocol)
	}
	var padLength uint8
	if fh.Flags.Has(FlagPushPromisePadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, err
		}
	}
	p, pp.PromiseID, err = readUint32(p)
	if err != nil {
		return nil, err
	}
	pp.PromiseID = pp.PromiseID & PadBit
	var payloadLength = len(p)
	if int(padLength) > payloadLength {
		return nil, ConnectionError(ErrCodeProtocol)
	}
	pp.HeaderFragBuf = p[:payloadLength-int(padLength)]
	return pp, nil
}

