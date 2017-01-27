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
				dy.heartBeat(time.Now().Unix())
				log.Println("heart beat")
			}
		}
	}()
}

func (dy *Douyu) heartBeat(now int64) {
	s := bytes.Join([][]byte{[]byte("type@=keeplive/tick@="), number2bytes(now)}, []byte(""))
	s = PackRequest(s)
	dy.Write(s)
}
