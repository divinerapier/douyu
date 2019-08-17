package douyu

import (
	"fmt"
	"strings"
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
		message := <-dy.chatMsgChan
		user := dy.getMessageField(message, "nn")
		text := dy.getMessageField(message, "txt")
		fmt.Printf("%v %20s: %s\n", time.Now(), user, text)
	}
}

func (dy *Douyu) getMessageField(message Message, field string) string {
	// eg: message is type@=lgpoolsite/zone@=1/deadsec@=17070/
	// we want get field zone
	if !strings.HasSuffix(field, "@=") {
		field = field + "@="
	}
	index := bytes.Index(message, []byte(field))
	if index == -1 {
		return ""
	}
	// zone@=1/deadsec@=17070
	message = message[index:]
	index = bytes.IndexByte(message, '/')
	if index == -1 {
		return string(message[len(field):])
	}
	return string(message[len(field):index])
}
