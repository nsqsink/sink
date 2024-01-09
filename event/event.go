package event

import "github.com/nsqsink/sink/contract"

type Module struct {
	topic         string   // topic name
	sourceAddress []string // source of the topic, for nsq its a nsqlookupd address
}

// NewEvent create new event
func NewEvent(topic string, source []string) contract.Event {
	return Module{topic: topic, sourceAddress: source}
}
