package http2

// HeadersFrame headers frame
type HeadersFrame struct {
	FrameHeader

	// Priority is set if FlagHeadersPriority is set in the FrameHeader.
	Priority PriorityParam

	HeaderFragBuf []byte // not owned
}

// HeadersEnded 头是否为分片
func (f *HeadersFrame) HeadersEnded() bool {
	return f.FrameHeader.Flags.Has(FlagHeadersEndHeaders)
}

// StreamEnded 流是否结束
func (f *HeadersFrame) StreamEnded() bool {
	return f.FrameHeader.Flags.Has(FlagHeadersEndStream)
}

// HasPriority 是否有 priority
func (f *HeadersFrame) HasPriority() bool {
	return f.FrameHeader.Flags.Has(FlagHeadersPriority)
}

// ParseHeadersFrame parse headers frame
func ParseHeadersFrame(fh *FrameHeader, p []byte) (*HeadersFrame, error) {
	var err error
	hf := &HeadersFrame{
		FrameHeader: *fh,
	}
	if fh.StreamID == 0 {
		return nil, ConnError{ErrCodeProtocol, "HEADERS frame with stream ID 0"}
	}
	var padLength uint8
	if fh.Flags.Has(FlagHeadersPadded) {
		if p, padLength, err = readByte(p); err != nil {
			return nil, err
		}
	}
	if fh.Flags.Has(FlagHeadersPriority) {
		var v uint32
		p, v, err = readUint32(p)
		if err != nil {
			return nil, err
		}
		hf.Priority.StreamDep = v & PadBit
		hf.Priority.Exclusive = (v != hf.Priority.StreamDep) // high bit was set
		p, hf.Priority.Weight, err = readByte(p)
		if err != nil {
			return nil, err
		}
	}
	if len(p)-int(padLength) <= 0 {
		return nil, streamError(fh.StreamID, ErrCodeProtocol)
	}
	hf.HeaderFragBuf = p[:len(p)-int(padLength)]
	return hf, nil
}

