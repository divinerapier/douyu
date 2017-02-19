package douyu

import (
	"bytes"
	"io"
	"os"

	log "qiniupkg.com/x/log.v7"
)

var (
	// TypeLoginRes 登陆响应类型
	TypeLoginRes = []byte("type@=loginres")
	// TypeKeepLive 心跳响应类型
	TypeKeepLive = []byte("type@=keeplive")
	// TypeChatmsg 弹幕消息
	TypeChatmsg = []byte("type@=chatmsg")
)

func (dy *Douyu) PrintResponse() {

	var rawMessageQueue = make(chan []byte, 1024)
	responseMsg := processResponseMessage(rawMessageQueue)

	go func() {
		for {
			select {
			case m := <-responseMsg:
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

func processResponseMessage(input <-chan []byte) <-chan *DouyuChatMessage {

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

// ReceiveResponse 接收回复消息
func (dy *Douyu) ReceiveResponse() {

	go dy.parseChatResponse()

	go func() {

		for {
			var buf [10240]byte
			cnt, err := dy.Read(buf[:])
			if err != nil {
				log.Error("receive response:", err)
				continue
			}
			if bytes.Contains(buf[:cnt], TypeChatmsg) {
				dy.chatMsgChan <- buf[:cnt]
			} else if bytes.Contains(buf[:cnt], TypeKeepLive) {
				dy.keepLiveChan <- buf[:cnt]
			} else {
				// log.Errorf("unknown type: [%s]\n", buf[12:cnt])
			}
		}
	}()
}
