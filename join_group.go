package douyu

import (
	"bytes"

	log "qiniupkg.com/x/log.v7"
)

// JoinGroupRequest 入组请求，固定为-9999
func (dy *Douyu) JoinGroupRequest(groupID int64) {

	joinGroupRequest(dy, groupID)
}

func joinGroupRequest(dy *Douyu, groupID int64) {
	data := bytes.Join([][]byte{[]byte("type@=joingroup/rid@="), number2bytes(dy.RoomID), []byte("/gid@=-9999/")}, []byte(""))

	data = PackRequest(data)
	if _, err := dy.Write(data); err != nil {
		log.Error("join_group_request send data err: ", err)
		return
	}

	// var buf [1024]byte
	// cnt, err := dy.Read(buf[:])
	// if err != nil {
	// 	log.Error("join_group_request recv data err: ", err)
	// 	return
	// }
	// log.Infof("dump join group response: %s", buf[:cnt])
}
