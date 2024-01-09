package nsq

import (
	"context"

	"github.com/nsqsink/sink/contract"
)

// NSQModule struct
// struct for
type NSQModule struct {
	consumers []contract.Consumer
}

// New
// return result initialization of NSQModule consumer
func New() contract.Streamer {
	module := &NSQModule{
		consumers: make([]contract.Consumer, 0),
	}

	return module
}

// RegisterConsumer implementation of register consumer method
func (m *NSQModule) RegisterConsumer(ctx context.Context, c contract.Consumer) error {
	// adding consumer to list of consumer on streamer
	m.consumers = append(m.consumers, c)
	return nil
}

// Run to run all handler in the consumer
func (m *NSQModule) Run() error {
	var err error

	// need to start all consumer
	for _, c := range m.consumers {
		if err = c.Run(); err != nil {
			return err
		}
	}

	return nil
}

// Stop to stop all consumer handler in the consumer
// using go routine to make it faster
func (m *NSQModule) Stop() error {
	errChan := make(chan error)

	for _, c := range m.consumers {
		go func(nsqConsumer contract.Consumer) {
			var err error
			defer func() {
				errChan <- err
			}()

			err = nsqConsumer.Stop()
		}(c)
	}

	close(errChan)

	var err error
	for tempErr := range errChan {
		if tempErr != nil {
			err = tempErr
		}
	}

	return err
}
