package washtub

import (
	"context"

	entities "github.com/nsqsink/sink/entity"
)

type Washtuber interface {
	Pulse(ctx context.Context, data entities.PulseRequest) chan error

	Message(ctx context.Context, data entities.MessageRequest) (*entities.MessageResponse, error)
}
