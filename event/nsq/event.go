package nsq

import "github.com/nsqsink/sink/contract"

type Module struct {
	topic         string   // topic name
	sourceAddress []string // source of the topic, for nsq its a nsqlookupd address
}

// NewEvent create new event
func New(topic string, source []string) contract.Event {
	return Module{topic: topic, sourceAddress: source}
}

// GetTopic return topic name
func (e Module) GetTopic() string {
	return e.topic
}

// GetSource return the source address for the topic
func (e Module) GetSource() []string {
	return e.sourceAddress
}
