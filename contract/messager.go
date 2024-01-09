package contract

import (
	"time"
)

// Messager
// contract for message data
type Messager interface {
	Finish()
	Requeue(delay time.Duration)
	// RequeueWithoutBackoff(delay time.Duration)
	GetAttempts() uint16
	GetBody() []byte
}
