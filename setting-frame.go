package http2

import "encoding/binary"

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
