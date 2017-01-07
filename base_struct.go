package douyu

import (
	"encoding/binary"
	"fmt"
	"net"
)

// OFFSET 数据长度与总长度的偏移量
const (
	Offset       = 4
	HeaderLength = 8
)

type Douyu struct {
	net.Conn
	*DouyuLoginResponse
	RoomID int64
}

type DouyuMessage struct {
	Length uint32
	Header *DouyuMessageHeader
	Data   []byte // '\0' 结尾
}

type DouyuMessageHeader struct {
	Length        uint32
	Type          uint16 // 689 cli -> srv      690   srv -> cli
	EncField      uint8  // 0
	ReservedField uint8  // 0
}

type DouyuLoginResponse struct {
	Type            []byte
	UserID          int64
	RoomGroup       int64
	Pg              int64
	SessionID       int64
	Username        []byte
	Nickname        []byte
	IsSigned        bool
	SignedCount     int64
	LiveStat        bool
	NeedPhoneVerify bool
	BestDlev        int64
	CurLev          int64
	ErrCode         int64
}

func PackRequest(data []byte) []byte {
	var dm DouyuMessage
	dm.Data = data
	dm.Data = append(dm.Data, byte(0))
	dm.Length = HeaderLength + 1 + uint32(len(data))
	dm.Header = &DouyuMessageHeader{
		Length: dm.Length,
		Type:   689,
	}
	return dm.Marshal()
}

func (dm *DouyuMessage) Marshal() []byte {
	buf := make([]byte, dm.Length+Offset)
	binary.LittleEndian.PutUint32(buf[:4], dm.Length)
	binary.LittleEndian.PutUint32(buf[4:8], dm.Header.Length)
	binary.LittleEndian.PutUint16(buf[8:10], dm.Header.Type)
	copy(buf[12:], dm.Data[0:])
	return buf
}

func (dy *Douyu) String() string {
	var s = `{"local":%s, "remote":%s, "room_id":%d, "info":%s}`
	return fmt.Sprintf(s, dy.LocalAddr(), dy.RemoteAddr(), dy.RoomID, dy.DouyuLoginResponse)
}

func (lr *DouyuLoginResponse) String() string {
	s := `{"type":%s,"user_id":%d,"room_group":%d,"pg":%d,"session_id":%d,"username":%s,"nickname":%s,"is_signed":%v,"signed_count":%d,"live_stat":%v,"need_phone_verify":%v,"best_dlev":%d,"cur_lev":%d}`
	return fmt.Sprintf(s, lr.Type, lr.UserID, lr.RoomGroup, lr.Pg, lr.SessionID, lr.Username, lr.Nickname, lr.IsSigned, lr.SignedCount, lr.LiveStat, lr.NeedPhoneVerify, lr.BestDlev, lr.CurLev)
}

func (dym *DouyuMessage) String() string {
	s := `{"length":%d,"header":%s,"data":%s}`
	return fmt.Sprintf(s, dym.Length, dym.Header, dym.Data)
}

func (dymh *DouyuMessageHeader) String() string {
	s := `{"length":%d,"type":%d,"enc_field":%d,"reserved_field":%d}`
	return fmt.Sprintf(s, dymh.Length, dymh.Type, dymh.EncField, dymh.ReservedField)
}
