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
				log.Println(m)
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

	output := make(chan *DouyuChatMessage, 1024)
	go func() {

		for {
			select {
			case msg := <-input:

				chatMsg := decodeChatMessage(msg)
				if chatMsg == nil {
					continue
				}

				output <- chatMsg
			}
		}

	}()
	return output
}

func decodeChatMessage(msg []byte) *DouyuChatMessage {
	if len(msg) == 0 {
		return nil
	}

	chatType := []byte("type@=chatmsg")
	heartBeatType := []byte("type@=keeplive")

	if bytes.Contains(msg, heartBeatType) {
		log.Println(string(msg))
		return nil
	}

	nn := []byte("/nn@=")
	txt := []byte("/txt@=")
	nickNameBegin, nickNameEnd, txtBegin, txtEnd := 0, 0, 0, 0

	chatMsg := AcquireChatmessage()
	chatMsg.Time = time.Now()
	if !bytes.Contains(msg, chatType) {
		return nil
	}
	chatMsg.Time = time.Now()
	if nickNameBegin = bytes.Index(msg, nn); nickNameBegin < 0 {
		ReleaseChatmessage(chatMsg)
		return nil
	} else {
		nickNameBegin += len(nn)
		nickNameEnd = nickNameBegin + bytes.IndexByte(msg[nickNameBegin:], '/')
		chatMsg.Username = msg[nickNameBegin:nickNameEnd]
	}
	if txtBegin = bytes.Index(msg, txt); txtBegin < 0 {
		ReleaseChatmessage(chatMsg)
		return nil
	} else {
		txtBegin += len(txt)
		txtEnd = txtBegin + bytes.IndexByte(msg[txtBegin:], '/')
		chatMsg.Message = msg[txtBegin:txtEnd]
	}

	return nil
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

func (dy *Douyu) parseChatResponse() {

	for {
		msg := <-dy.chatMsgChan
		log.Infof("dump chat message: %s\n", msg[12:])
		lines := bytes.Split(msg, []byte("/"))
		var rid, gid, uid, nn, txt string
		for _, v := range lines {
			kv := bytes.Split(v, []byte("@="))
			if len(kv) > 1 {
				switch string(kv[0]) {
				case "rid":
					rid = string(kv[1])
				case "gid":
					gid = string(kv[1])
				case "uid":
					uid = string(kv[1])
				case "nn":
					nn = string(kv[1])
				case "txt":
					txt = string(kv[1])
				default:
					// log.Error("unknown key:", string(kv[0]), "value:", string(kv[1]))
				}
			}
		}
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", rid, gid, uid, nn, txt)
	}
}
