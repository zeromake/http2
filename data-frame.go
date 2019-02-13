package http2

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
