package http2

// DataFrame body
type DataFrame struct {
	FrameHeader
	Data []byte
}

// ParseDataFrame 解析
func ParseDataFrame(fh *FrameHeader, payload []byte) (*DataFrame, error) {
	if fh.StreamID == 0 {
		return nil, ConnError{ErrCodeProtocol, "DATA frame with stream ID 0"}
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
		return nil, ConnError{ErrCodeProtocol, "pad size larger than data payload"}
	}
	f.Data = payload[:len(payload)-int(padSize)]
	return f, nil
}
