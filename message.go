package douyu

import "bytes"

var (
	// TypeLoginRes 登陆响应类型
	TypeLoginRes = []byte("type@=loginres")
	// TypeKeepLive 心跳响应类型
	TypeKeepLive = []byte("type@=keeplive")
	// TypeChatmsg 弹幕消息
	TypeChatmsg = []byte("type@=chatmsg")
	// Some magic message types
	TypeLgpoolsite  = []byte("type@=lgpoolsite")
	TypeNodeNumInfo = []byte("type@=noble_num_info")
	TypeShrn        = []byte("type@=shrn")
	TypeUenter      = []byte("type@=uenter")
	TypeUlRandlist  = []byte("type@=ul_ranklist")
)

type Message []byte

func (message Message) Type() string {
	switch {
	case bytes.HasPrefix(message, TypeLoginRes):
		return "loginres"
	case bytes.HasPrefix(message, TypeKeepLive):
		return "keeplive"
	case bytes.HasPrefix(message, TypeChatmsg):
		return "chatmsg"
	case bytes.HasPrefix(message, TypeLgpoolsite):
		return "lgpoolsite"
	case bytes.HasPrefix(message, TypeNodeNumInfo):
		return "noble_num_info"
	case bytes.HasPrefix(message, TypeShrn):
		return "shrn"
	case bytes.HasPrefix(message, TypeUenter):
		return "uenter"
	case bytes.HasPrefix(message, TypeUlRandlist):
		return "ul_ranklist"
	default:
		return "unknown type"
	}
}

func (message Message) IsChatMessage() bool {
	return bytes.HasPrefix(message, TypeChatmsg)
}

func (message Message) UnknownType() bool {
	return message.Type() == "unknown type"
}
