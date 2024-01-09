package nsq

import "github.com/nsqsink/sink/log"

// Config config for consumer
type Config struct {
	ChannelName string // name of the consumer channel
	Concurrent  int    // number of concurrent consumer
	MaxAttempt  int    // max attempt of consumer to handle a message
	MaxInFlight int
	LogLevel    log.LogLevel // setting for log level (1 - )
}
