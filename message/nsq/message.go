package nsq

import (
	"github.com/nsqio/go-nsq"
	"github.com/nsqsink/sink/contract"
)

type Message struct {
	*nsq.Message
}

// New
// return message object
func New(msg *nsq.Message) contract.Messager {
	return Message{msg}
}

// GetAttempts
// return number of attempts consume this message
func (m Message) GetAttempts() uint16 {
	return m.Attempts
}

// GetBody
// return body of the message in byte
func (m Message) GetBody() []byte {
	return m.Body
}
