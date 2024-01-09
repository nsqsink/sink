package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqsink/sink/config"
	consumer "github.com/nsqsink/sink/consumer/nsq"
	event "github.com/nsqsink/sink/event/nsq"
	"github.com/nsqsink/sink/handler"
	logger "github.com/nsqsink/sink/log"
	streamer "github.com/nsqsink/sink/streamer/nsq"
)

func main() {
	// init configuration variables
	var (
		cfg            config.App
		consumerCfg    config.Consumer
		configFilePath string
	)

	// parse configuration from flags
	// file configuration
	flag.StringVar(&configFilePath, "config-path", "", "Define config file path if you want to read config from files instead of flag")

	flag.StringVar(&cfg.LogLevel, "log-level", "", "Define log level")

	// single consumer config
	flag.StringVar(&consumerCfg.ID, "id", "", "Define consumer id")
	flag.StringVar(&consumerCfg.Topic, "topic", "", "Define the topic name")
	flag.StringVar(&consumerCfg.Source, "source", "", "Define the source of the topic")
	flag.IntVar(&consumerCfg.Concurrent, "concurrent", 1, "Define the number of concurrent for the consumer")
	flag.IntVar(&consumerCfg.MaxAttempt, "max-attempt", 5, "Define the number of max attempt for the consumer to process the message")
	flag.IntVar(&consumerCfg.MaxInFlight, "max-in-flight", 5, "Define the number of max in flight for the consumer to process the message")
	flag.BoolVar(&consumerCfg.Active, "active", false, "Define the consumer is active or not")
	flag.StringVar(&consumerCfg.Sinker.Type, "sinker-type", "", "Define the type of the sinker, example: http")
	flag.StringVar(&consumerCfg.Sinker.Parser.Template, "parser-template", "", "Define the template for the parser")
	flag.StringVar(&consumerCfg.Sinker.Parser.Type, "parser-type", "", "Define the type of the parser")

	flag.Parse()

	// load config from files
	if configFilePath == "" {
		cfg.Consumers = append(cfg.Consumers, consumerCfg)
	} else {
		jsonFile, err := os.Open(configFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Successfully open config file %s\n", configFilePath)
		defer jsonFile.Close()

		configByte, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Fatalln(err)
		}

		if err := json.Unmarshal(configByte, &cfg.Consumers); err != nil {
			log.Fatalln(err)
		}
	}

	// validate config
	if _, err := cfg.Validate(); err != nil {
		log.Fatalln(err)
	}

	// create streamer server
	streamerModule := streamer.New()

	// for each of consumer, register consumer
	for _, cfgConsumer := range cfg.Consumers {
		ctx := context.Background()

		// create event
		e := event.New(cfgConsumer.Topic, cfgConsumer.ParseSource())

		// create handler (still dummy)
		h, err := handler.New(cfgConsumer.Sinker)
		if err != nil {
			log.Fatalln(err)
		}

		// create consumer
		c, err := consumer.New(ctx, e, h, consumer.Config{
			ChannelName: cfgConsumer.ID,
			Concurrent:  cfgConsumer.Concurrent,
			MaxAttempt:  cfgConsumer.MaxAttempt,
			MaxInFlight: cfgConsumer.MaxInFlight,
			LogLevel:    logger.LogLevel(cfg.LogLevel),
		})
		if err != nil {
			log.Fatalln(err)
		}

		if err := streamerModule.RegisterConsumer(ctx, c); err != nil {
			log.Fatalln(err)
		}
	}

	// run using socketmaster
	idleConnsClosed := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)

		// When using socketmaster, it send SIGTERM after spawning new process,
		// as it mentioned on the doc https://github.com/tokopedia/socketmaster#how-it-works
		// SIGHUP is for handling upstart reload
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
		sign := <-signals

		log.Printf("signal %s terminate detected, gracefully shutdown the workers", sign)

		// We received an os signal, shut down.
		if err := streamerModule.Stop(); err != nil {
			// Error from closing listeners, or context timeout:
			fmt.Printf("server Shutdown: %v\n", err)
		}

		close(idleConnsClosed)
	}()

	fmt.Println("server running")
	if err := streamerModule.Run(); err != nil {
		// Error starting or closing consumer:
		log.Fatalln(err)
	}

	<-idleConnsClosed
	fmt.Println("server shutdown gracefully")
}
