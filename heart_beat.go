package douyu

import (
	"bytes"
	"time"
)

// HeartBeat 心跳消息 每 45 s
func (dy *Douyu) HeartBeat() {
	tick := time.Tick(time.Second * 15)
	go func() {
		for {
			select {
			case <-tick:
				dy.heartBeat()
			}
		}
	}()
}

func (dy *Douyu) heartBeat() {
	s := bytes.Join([][]byte{}, []byte(""))
	s = PackRequest(s)
	dy.Write(s)
}
