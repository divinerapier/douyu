package douyu

import (
	"bytes"

	log "qiniupkg.com/x/log.v7"
)

func (dy *Douyu) JoinGroupRequest(groupID int64) {

	data := bytes.Join([][]byte{[]byte("type@=joingroup/rid@="), number2bytes(dy.RoomID), []byte("/gid@=-9999/")}, []byte(""))

	data = PackRequest(data)
	if _, err := dy.Write(data); err != nil {
		log.Error("join_group_request send data err: ", err)
		return
	}

	var buf [1024]byte
	_, err := dy.Read(buf[:])
	if err != nil {
		log.Error("join_group_request recv data err: ", err)
		return
	}
}
