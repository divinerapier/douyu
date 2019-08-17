package danmaku

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrWrongMessageType = errors.New("wrong message type")
)

type Protocol struct {
	transport *Transport
}

type Header struct {
	length      uint32
	messageType uint16
}

type ChatMessage struct {
	Time    time.Time
	User    string
	Message string
}

func NewProtocol(trans *Transport) *Protocol {
	return &Protocol{
		transport: trans,
	}
}

func (h *Header) bodyLength() uint32 {
	return h.length - 8
}

func (p *Protocol) readHeader() (*Header, error) {
	var header [12]byte
	if err := p.transport.ReadFull(header[:]); err != nil {
		return nil, err
	}

	length := binary.LittleEndian.Uint32(header[:4])
	length2 := binary.LittleEndian.Uint32(header[4:8])
	messageType := binary.LittleEndian.Uint16(header[8:10])
	if length != length2 || length <= 0 || messageType != 690 {
		return nil, fmt.Errorf("corrupted header received. len: %d, len2: %d, message type: %d",
			length, length2, messageType)
	}

	return &Header{
		length:      length,
		messageType: 690,
	}, nil
}

func (p *Protocol) readBody(length uint32) ([]byte, error) {
	return p.transport.ReadSize(int(length))
}

func (p *Protocol) Read() ([]byte, error) {

	// read header
	header, err := p.readHeader()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read header")
	}
	// read body

	body, err := p.readBody(header.bodyLength())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}
	return body, nil
}

func (p *Protocol) Write(message []byte) error {
	length := len(message)
	data := make([]byte, 13+length)
	binary.LittleEndian.PutUint32(data[:4], uint32(length+9))
	binary.LittleEndian.PutUint32(data[4:8], uint32(length+9))
	binary.LittleEndian.PutUint16(data[8:10], 689)
	copy(data[12:], []byte(message))

	err := p.transport.Write(data)
	if err != nil {
		return errors.Wrapf(err, "write message. '%s'", message)
	}
	return nil
}

func (p *Protocol) WriteString(message string) error {
	return p.Write([]byte(message))
}

func (p *Protocol) ReadChatMessage() (*ChatMessage, error) {
	message, err := p.Read()
	if err != nil {
		return nil, err
	}
	var chatMessage ChatMessage
	err = NewChatMessageDecoder().Decode(message, &chatMessage)
	return &chatMessage, err
}

type ChatMessageDecoder struct{}

func NewChatMessageDecoder() ChatMessageDecoder {
	return ChatMessageDecoder{}
}

func (decoder ChatMessageDecoder) Decode(message []byte, v interface{}) error {
	chatMessage, ok := v.(*ChatMessage)
	if !ok {
		return fmt.Errorf("wrong target type. expect *ChatMessage, get %T", v)
	}

	if len(message) == 0 {
		return errors.New("invalid message")
	}

	if msgType := decoder.getMessageField(message, "type"); msgType != "chatmsg" {
		return ErrWrongMessageType
	}

	nn := decoder.getMessageField(message, "nn")
	txt := decoder.getMessageField(message, "txt")
	if nn == "" || txt == "" {
		return fmt.Errorf("invalid chat message. %s", message)
	}

	chatMessage.Time = time.Now()
	chatMessage.User = nn
	chatMessage.Message = txt
	return nil
}

func (dy *ChatMessageDecoder) getMessageField(message []byte, field string) string {
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
