package http2

import (
	"fmt"
)

// ErrCode 错误定义
type ErrCode uint32

const (
	// ErrCodeNo 未知错误
	ErrCodeNo ErrCode = 0x0
	// ErrCodeProtocol 协议错误
	ErrCodeProtocol ErrCode = 0x1
	// ErrCodeInternal 网络错误
	ErrCodeInternal ErrCode = 0x2
	// ErrCodeFlowControl 流量控制错误
	ErrCodeFlowControl ErrCode = 0x3
	// ErrCodeSettingsTimeout 等待 Settings Frame 超时
	ErrCodeSettingsTimeout ErrCode = 0x4
	// ErrCodeStreamClosed 流关闭错误
	ErrCodeStreamClosed ErrCode = 0x5
	// ErrCodeFrameSize frame 大小超过限制
	ErrCodeFrameSize ErrCode = 0x6
	// ErrCodeRefusedStream 流拒绝错误
	ErrCodeRefusedStream ErrCode = 0x7
	// ErrCodeCancel 取消错误
	ErrCodeCancel ErrCode = 0x8
	// ErrCodeCompression 压缩错误
	ErrCodeCompression ErrCode = 0x9
	// ErrCodeConnect 连接错误
	ErrCodeConnect            ErrCode = 0xa
	ErrCodeEnhanceYourCalm    ErrCode = 0xb
	ErrCodeInadequateSecurity ErrCode = 0xc
	ErrCodeHTTP11Required     ErrCode = 0xd
)

var errCodeName = map[ErrCode]string{
	ErrCodeNo:                 "NO_ERROR",
	ErrCodeProtocol:           "PROTOCOL_ERROR",
	ErrCodeInternal:           "INTERNAL_ERROR",
	ErrCodeFlowControl:        "FLOW_CONTROL_ERROR",
	ErrCodeSettingsTimeout:    "SETTINGS_TIMEOUT",
	ErrCodeStreamClosed:       "STREAM_CLOSED",
	ErrCodeFrameSize:          "FRAME_SIZE_ERROR",
	ErrCodeRefusedStream:      "REFUSED_STREAM",
	ErrCodeCancel:             "CANCEL",
	ErrCodeCompression:        "COMPRESSION_ERROR",
	ErrCodeConnect:            "CONNECT_ERROR",
	ErrCodeEnhanceYourCalm:    "ENHANCE_YOUR_CALM",
	ErrCodeInadequateSecurity: "INADEQUATE_SECURITY",
	ErrCodeHTTP11Required:     "HTTP_1_1_REQUIRED",
}

func (e ErrCode) String() string {
	if s, ok := errCodeName[e]; ok {
		return s
	}
	return fmt.Sprintf("unknown error code 0x%x", uint32(e))
}

// ConnectionError 连接错误
type ConnectionError ErrCode

func (e ConnectionError) Error() string { return fmt.Sprintf("connection error: %s", ErrCode(e)) }

// StreamError 流错误
type StreamError struct {
	StreamID uint32
	Code     ErrCode
	Cause    error
}

func streamError(id uint32, code ErrCode) StreamError {
	return StreamError{StreamID: id, Code: code}
}

func (e StreamError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("stream error: stream ID %d; %v; %v", e.StreamID, e.Code, e.Cause)
	}
	return fmt.Sprintf("stream error: stream ID %d; %v", e.StreamID, e.Code)
}

type connError struct {
	Code   ErrCode
	Reason string
}

func (e connError) Error() string {
	return fmt.Sprintf("http2: connection error: %v: %v", e.Code, e.Reason)
}
