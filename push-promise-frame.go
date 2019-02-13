package http2

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
