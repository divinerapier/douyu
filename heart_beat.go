package douyu

import (
	"bytes"
	"time"

	log "qiniupkg.com/x/log.v7"
)

// HeartBeat 心跳消息 每 45 s
func (dy *Douyu) HeartBeat() {
	tick := time.Tick(time.Second * 45)
	go func() {
		for {
			select {
			case <-tick:
				now := time.Now().Unix()
				resp := dy.heartBeat(now)
				log.Info("heart beat time:", now, "response:", resp)
			}
		}

		// for {
		// 	now := time.Now().Unix()
		// 	resp := dy.heartBeat(now)
		// 	log.Info("heart beat time:", now, "response:", resp)
		// 	time.Sleep(time.Second * 15)
		// }

	}()
}

func (dy *Douyu) heartBeat(now int64) int64 {
	s := bytes.Join([][]byte{[]byte("type@=keeplive/tick@="), number2bytes(now), []byte{'/'}}, []byte(""))
	s = PackPacket(s)
	log.Info("heart beat:", s)
	dy.Write(s)
	msg := <-dy.keepLiveChan
	log.Info("dump keeplive msg:", string(msg[12:]))
	start := bytes.Index(msg, []byte("tick@="))
	if start < 0 {
		return -1
	}
	start += 6
	end := bytes.IndexByte(msg[start:], '/')
	if end < 0 {
		log.Error("keep live msg error: end of tick not found. msg:", string(msg[12:]))
	}
	data := msg[start : start+end]

	return bytes2number(data)

	// return now
}
