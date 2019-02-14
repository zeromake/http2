package http2

// Flags flag type
type Flags uint8

// Has Flags has a flag
func (f Flags) Has(v Flags) bool {
	return (f & v) == v
}

const (
	// Data Frame

	// FlagDataEndStream end stream
	FlagDataEndStream Flags = 0x1
	// FlagDataPadded data padded
	FlagDataPadded Flags = 0x8

	// Headers Frame

	// FlagHeadersEndStream end stream
	FlagHeadersEndStream Flags = 0x1
	// FlagHeadersEndHeaders end headers
	FlagHeadersEndHeaders Flags = 0x4
	// FlagHeadersPadded header has padded
	FlagHeadersPadded Flags = 0x8
	// FlagHeadersPriority header has prority
	FlagHeadersPriority Flags = 0x20

	// Settings Frame

	// FlagSettingsAck settings frame is ack
	FlagSettingsAck Flags = 0x1

	// Ping Frame

	// FlagPingAck ping frame is ack
	FlagPingAck Flags = 0x1

	// Continuation Frame

	// FlagContinuationEndHeaders contnuation frame is headers end
	FlagContinuationEndHeaders Flags = 0x4

	// PushPromise Frame

	// FlagPushPromiseEndHeaders PushPromise frame is header end
	FlagPushPromiseEndHeaders Flags = 0x4
	// FlagPushPromisePadded PushPromise frame has padded
	FlagPushPromisePadded Flags = 0x8
)
