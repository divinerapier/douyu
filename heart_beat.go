package douyu

import (
	"bytes"
	"time"

	log "qiniupkg.com/x/log.v7"
)

// HeartBeat 心跳消息 每 45 s
func (dy *Douyu) HeartBeat() {
	tick := time.Tick(time.Second * 15)
	go func() {
		for {
			select {
			case <-tick:
				now := time.Now().Unix()
				resp := dy.heartBeat(now)
				log.Println("heart beat time:", now, "response:", resp)
			}
		}
	}()
}

func (dy *Douyu) heartBeat(now int64) int64 {
	s := bytes.Join([][]byte{[]byte("type@=keeplive/tick@="), number2bytes(now), []byte{'/'}}, []byte(""))
	s = PackRequest(s)
	dy.Write(s)
	msg := <-dy.keepLiveChan
	start := bytes.Index(msg, []byte("tick@="))
	if start < 0 {
		return -1
	}
	start += 6
	end := bytes.IndexByte(msg[start:], '/')
	if end < 0 {
		log.Error("keep live msg error: end of tick not found. msg:", string(msg))
	}
	data := msg[start : start+end]

	return bytes2number(data)
}
