package douyu

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	log "qiniupkg.com/x/log.v7"
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
			message, err := dy.readMeaage()
			if err != nil {
				log.Error("receive response:", err)
				if err == io.EOF {
					log.Fatal("read eof from remote. exit -1")
					return
				}
				continue
			}
			if message.IsChatMessage() {

			} else if message.UnknownType() {
				fmt.Fprintf(os.Stderr, "UNKNOWN MESSAGE TYPE. '%s'", message)
			} else {
				// do nothing
			}
		}
	}()
}

func (dy *Douyu) readMeaage() (Message, error) {
	var header [12]byte
	// read header
	_, err := io.ReadFull(dy, header[:])
	if err != nil {
		panic(err)
	}
	length := binary.LittleEndian.Uint32(header[:4])
	length2 := binary.LittleEndian.Uint32(header[4:8])
	messageType := binary.LittleEndian.Uint32(header[8:10])
	if length != length2 || length <= 0 || messageType != 690 {
		panic(fmt.Errorf("corrupted header received. len: %d, len2: %d, message type: %d",
			length, length2, messageType))
	}
	// read body
	body := make([]byte, length-8)
	_, err = io.ReadFull(dy, body)
	if err != nil {
		panic(err)
	}
	return Message(body), nil
}
