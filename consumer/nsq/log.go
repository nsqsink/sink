package nsq

import (
	"github.com/nsqio/go-nsq"
	"github.com/nsqsink/sink/log"
)

// toNSQLogLevel return log level for NSQ
func toNSQLogLevel(l log.LogLevel) nsq.LogLevel {
	switch l {
	case "debug":
		return nsq.LogLevelDebug
	case "info":
		return nsq.LogLevelInfo
	case "warn":
		return nsq.LogLevelWarning
	case "error":
		return nsq.LogLevelError
	}

	return nsq.LogLevelMax
}
