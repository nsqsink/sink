package nsq

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
	"github.com/nsqsink/sink/contract"
	message "github.com/nsqsink/sink/message/nsq"
)

type Module struct {
	nsqConsumer *nsq.Consumer
	source      []string
}

// New return consumer module / object
// accepting event of the message, the handler for the event and the configuration of the consumer
func New(ctx context.Context, e contract.Event, h contract.Handler, cfg Config) (contract.Consumer, error) {
	// checking required data
	if e.GetTopic() == "" {
		return nil, errors.New("empty event topic")
	}

	if len(e.GetSource()) == 0 {
		return nil, errors.New("empty event source")
	}

	if h == nil {
		return nil, errors.New("empty handler")
	}

	// checking max attempt
	// using default value if empty
	if cfg.MaxAttempt <= 0 {
		cfg.MaxAttempt = constDefaultMaxAttempt
	}

	// checking max in flight
	// using default value if empty
	if cfg.MaxInFlight <= 0 {
		cfg.MaxInFlight = constDefaultMaxInflight
	}

	// generate random channel name from uuid if empty
	if cfg.ChannelName == "" {
		cfg.ChannelName = e.GetTopic() + "-" + uuid.NewString()
	}

	// create new consumer
	c, err := nsq.NewConsumer(e.GetTopic(), cfg.ChannelName, &nsq.Config{
		MaxAttempts: uint16(cfg.MaxAttempt),
		MaxInFlight: cfg.MaxInFlight,
	})
	if err != nil {
		return nil, err
	}

	// set log level
	c.SetLoggerLevel(toNSQLogLevel(cfg.LogLevel))

	// wrap handler to nsq handler
	handlerFn := func(msg *nsq.Message) error {
		return h.Handle(message.New(msg))
	}

	// add handler based on concurrent numbers
	if cfg.Concurrent > 0 {
		c.AddConcurrentHandlers(nsq.HandlerFunc(handlerFn), cfg.Concurrent)
	} else {
		c.AddHandler(nsq.HandlerFunc(handlerFn))
	}

	// return consumer
	return Module{
		nsqConsumer: c,
		source:      e.GetSource(),
	}, nil
}

// Run is a method to run / start the consumer to listen from an event
func (m Module) Run() error {
	// run the consumer by connecting to nsqlookupd
	if err := m.nsqConsumer.ConnectToNSQDs(m.source); err != nil {
		return err
	}

	return nil
}

// Stop is a method to stop and close the consumer from listening an event
func (m Module) Stop() error {
	m.nsqConsumer.Stop()
	return nil
}
