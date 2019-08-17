package danmaku

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const DouyuDanmakuServer = "openbarrage.douyutv.com:8601"

// Danmaku douyu danmu client
type Danmaku struct {
	protocol *Protocol
	roomID   int64
}

// Dial 打开斗鱼弹幕
func Dial(address string, opts ...DialOption) (*Danmaku, error) {
	var danmaku Danmaku

	for _, opt := range opts {
		opt(&danmaku)
	}
	transport, err := NewTransport("tcp", address)
	if err != nil {
		return nil, err
	}
	danmaku.protocol = NewProtocol(transport)
	return &danmaku, nil
}

type DialOption func(*Danmaku)

func WithRoom(room int64) DialOption {
	return func(danmaku *Danmaku) {
		danmaku.roomID = room
	}
}

func (danmaku *Danmaku) Run() error {

	if err := danmaku.login(); err != nil {
		return errors.Wrap(err, "failed to login")
	}

	if err := danmaku.joingroup(); err != nil {
		return errors.Wrap(err, "failed to join group")
	}

	go danmaku.keepalive()

	danmaku.readLoop()
	return nil
}

func (danmaku *Danmaku) keepalive() {
	tick := time.Tick(time.Second * 45)
	for {
		select {
		case <-tick:
			if err := danmaku.protocol.WriteString("type@=mrkl/"); err != nil {
				logrus.Errorf("keepalive error: %v", err)
			}
		}
	}
}

func (danmaku *Danmaku) login() error {
	message := fmt.Sprintf("type@=loginreq/roomid@=%d", danmaku.roomID)
	return danmaku.protocol.WriteString(message)
}

func (danmaku *Danmaku) joingroup() error {
	data := fmt.Sprintf("type@=joingroup/rid@=%d/gid@=-9999/", danmaku.roomID)
	return danmaku.protocol.WriteString(data)
}

func (danmaku *Danmaku) readLoop() {
	for {
		message, err := danmaku.protocol.ReadChatMessage()
		if err == ErrWrongMessageType {
			continue
		}
		if err != nil {
			logrus.Errorf("failed to read message. error: %v", err)
			continue
		}
		logrus.Printf("%v  %s: %s\n", message.Time.Format("2006-01-02 15:04:05"), message.User, message.Message)
	}
}
