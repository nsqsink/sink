package contract

import "context"

// Sinker is the interface that wraps the basic Write and Close methods.
type Sinker interface {
	// Write writes the given data to the sink.
	Write(ctx context.Context, data []byte) ([]byte, error)

	// Close closes the sink.
	Close(ctx context.Context) error
}
