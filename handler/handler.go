package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nsqsink/sink/config"
	"github.com/nsqsink/sink/contract"
	jsonParser "github.com/nsqsink/sink/parser/json"
	"github.com/nsqsink/sink/sink/http"
)

// @TODO: still on progress
type Module struct {
	sinker contract.Sinker
	parser contract.Parser
}

func New(cfgSink config.Sinker) (contract.Handler, error) {
	var (
		sinker    contract.Sinker
		parser    contract.Parser
		errSinker error
		errParser error
	)

	// init sinker
	switch strings.ToLower(cfgSink.Type) {
	case "http":
		httpMethod := strings.ToUpper(cfgSink.HTTP.Method)
		sinker, errSinker = http.NewSink(cfgSink.HTTP.URL, httpMethod)
	default:
		errSinker = fmt.Errorf("sinker type %s not supported yet", cfgSink.Type)
	}
	if errSinker != nil {
		return nil, errSinker
	}

	// init parser
	switch strings.ToLower(cfgSink.Parser.Type) {
	case "json":
	default:
		parser, errParser = jsonParser.New(cfgSink.Parser)
	}
	if errParser != nil {
		return nil, errParser
	}

	return Module{
		sinker: sinker,
		parser: parser,
	}, nil
}

func (m Module) Handle(msg contract.Messager) error {
	ctx := context.Background()

	// parse body
	bodyMessage := msg.GetBody()

	log.Println(string(bodyMessage))

	// parse template
	parsed, err := m.parser.Parse(bodyMessage)
	if err != nil {
		log.Println(err)
		return err
	}

	// send to sinker
	_, err = m.sinker.Write(ctx, parsed)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
