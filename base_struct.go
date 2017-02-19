package douyu

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// OFFSET 数据长度与总长度的偏移量
const (
	Offset       = 4
	HeaderLength = 8
)

// Douyu douyu danmu client
type Douyu struct {
	net.Conn
	keepLiveChan chan []byte
	chatMsgChan  chan []byte
	*LoginResponse
	RoomID int64
}

// Message send to douyu server
type Message struct {
	Length uint32
	Header *MessageHeader
	Data   []byte // '\0' 结尾
}

// MessageHeader header of douyu message
type MessageHeader struct {
	Length        uint32
	Type          uint16 // 689 cli -> srv      690   srv -> cli
	EncField      uint8  // 0
	ReservedField uint8  // 0
}

// LoginResponse response struct for login request
type LoginResponse struct {
	Type            []byte `json:"type,omitempty,string"`
	UserID          int64  `json:"user_id,omitempty"`
	RoomGroup       int64  `json:"room_group,omitempty"`
	Pg              int64  `json:"pg,omitempty"`
	SessionID       int64  `json:"session_id,omitempty"`
	Username        []byte `json:"user_name,omitempty,string"`
	Nickname        []byte `json:"nick_name,omitempty,string"`
	IsSigned        bool   `json:"is_signed,omitempty"`
	SignedCount     int64  `json:"signed_count,omitempty"`
	LiveStat        bool   `json:"live_stat,omitempty"`
	NeedPhoneVerify bool   `json:"need_phone_verify,omitempty"`
	BestDlev        int64  `json:"best_delv,omitempty"`
	CurLev          int64  `json:"cur_lev,omitempty"`
	ErrCode         int64  `json:"err_code,omitempty"`
	IsIllegal       bool   `json:"is_illegal,omitempty"`
	IllCt           int64  `json:"ill_ct,omitempty"`
	IllTs           int64  `json:"ill_ts,omitempty"`
	Now             int64  `json:"now,omitempty"`
	Ps              int64  `json:"ps,omitempty"`
	Es              int64  `json:"es,omitempty"`
	It              int64  `json:"it,omitempty"`
	Its             int64  `json:"its,omitempty"`
	Nrc             int64  `json:"nrc,omitempty"`
	Ih              int64  `json:"ih,omitempty"`
	SID             int64  `json:"sid,omitempty"`
}

func PackRequest(data []byte) []byte {
	var dm Message
	dm.Data = data
	dm.Data = append(dm.Data, byte(0))
	dm.Length = HeaderLength + 1 + uint32(len(data))
	dm.Header = &MessageHeader{
		Length: dm.Length,
		Type:   689,
	}
	return dm.Marshal()
}

func (dy *Douyu) String() string {
	var s = `{"local":%s, "remote":"%s", "room_id":%d, "info":%s}`
	return fmt.Sprintf(s, dy.LocalAddr(), dy.RemoteAddr(), dy.RoomID, dy.LoginResponse)
}

func (lr *LoginResponse) String() string {
	data, _ := json.Marshal(lr)
	return string(data)
}

func (dm *Message) Marshal() []byte {
	buf := make([]byte, dm.Length+Offset)
	binary.LittleEndian.PutUint32(buf[:4], dm.Length)
	binary.LittleEndian.PutUint32(buf[4:8], dm.Header.Length)
	binary.LittleEndian.PutUint16(buf[8:10], dm.Header.Type)
	copy(buf[12:], dm.Data[0:])
	return buf
}

func (dym *Message) String() string {
	s := `{"length":%d,"header":%s,"data":"%s"}`
	return fmt.Sprintf(s, dym.Length, dym.Header, dym.Data)
}

func (dymh *MessageHeader) String() string {
	s := `{"length":%d,"type":%d,"enc_field":%d,"reserved_field":%d}`
	return fmt.Sprintf(s, dymh.Length, dymh.Type, dymh.EncField, dymh.ReservedField)
}
