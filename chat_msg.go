package douyu

import (
	"fmt"
	"time"

	"sync"

	"bytes"

	"io"

	"os"

	log "qiniupkg.com/x/log.v7"
)

func (dy *Douyu) ShowChatmessage() {
	var rawMessageQueue = make(chan []byte, 1024)
	chatmsg := processChatmessage(rawMessageQueue)

	go func() {
		for {
			select {
			case m := <-chatmsg:
				fmt.Println(m)
				ReleaseChatmessage(m)
			default:
				var buf [4096]byte
				cnt, err := dy.Read(buf[:])
				if err != nil {
					if err != io.EOF {
						log.Error("recv chat message err: ", err)
						continue
					} else {
						dy.Close()
						os.Exit(-1)
					}
				}
				// fmt.Printf("\n%s\n", buf[12:cnt])
				rawMessageQueue <- buf[12:cnt]
			}

		}
	}()
}

func processChatmessage(input <-chan []byte) <-chan *DouyuChatMessage {
	chatType := []byte("type@=chatmsg")
	nn := []byte("/nn@=")
	txt := []byte("/txt@=")
	output := make(chan *DouyuChatMessage, 1024)
	go func(chan *DouyuChatMessage) {
		for {
			select {
			case msg := <-input:
				if !bytes.Contains(msg, chatType) {
					continue
				}
				chatMsg := AcquireChatmessage()
				chatMsg.Time = time.Now()
				if begin := bytes.Index(msg, nn); begin < 0 {
					ReleaseChatmessage(chatMsg)
					continue
				} else {
					end := bytes.IndexByte(msg[begin+len(nn):], '/')
					chatMsg.Username = msg[begin+len(nn) : begin+len(nn)+end]
				}
				if begin := bytes.Index(msg, txt); begin < 0 {
					ReleaseChatmessage(chatMsg)
					continue
				} else {
					end := bytes.IndexByte(msg[begin+len(txt):], '/')
					chatMsg.Message = msg[begin+len(txt) : begin+len(txt)+end]
				}
				output <- chatMsg
			}
		}

	}(output)
	return output
}

type DouyuChatMessage struct {
	Time     time.Time
	Username []byte
	Message  []byte
}

var chatMessageFormatStr = "%d:%d:%d.%d\t%s\t\t\t%s"

func (dcm *DouyuChatMessage) String() string {
	now := dcm.Time
	return fmt.Sprintf(chatMessageFormatStr, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), dcm.Username, dcm.Message)
}

var douyuChatMessagePool = &sync.Pool{
	New: func() interface{} {
		return new(DouyuChatMessage)
	},
}

func AcquireChatmessage() *DouyuChatMessage {
	return douyuChatMessagePool.Get().(*DouyuChatMessage)
}

func ReleaseChatmessage(a *DouyuChatMessage) {
	if a != nil {
		ResetChatmessage(a)
		douyuChatMessagePool.Put(a)
	}
}

func ResetChatmessage(a *DouyuChatMessage) {
	a.Message = a.Message[:0]
	a.Username = a.Username[:0]
}
