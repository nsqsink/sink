package config

type (
	App struct {
		LogLevel  string     `json:"log_level"`
		Consumers []Consumer `json:"consumers"`
		Washtub   string     `json:"washtub"`
	}

	Consumer struct {
		ID          string `json:"id"`     // will be channel name
		Topic       string `json:"topic"`  // source topic
		Source      Source `json:"source"` // list source
		MaxAttempt  int    `json:"max_attempt"`
		MaxInFlight int    `json:"max_in_flight"`
		Concurrent  int    `json:"concurrent"`
		Sinker      Sinker `json:"sinker"`
		Active      bool   `json:"active"`
	}

	Source struct {
		NSQD       []string `json:"nsqd"`       // list source nsqd, separated by comma for multiple value
		NSQLookupd []string `json:"nsqlookupd"` // list source nsqd, separated by comma for multiple value
	}

	Sinker struct {
		Type   string       `json:"type"` //sinker type
		Parser Parser       `json:"parser"`
		Config SinkerConfig `json:"config"`
	}

	SinkerConfig struct {
		HTTP HTTPSinker `json:"http"`
		File FileSinker `json:"file"`
	}

	Parser struct {
		Type     string `json:"type"`     // json, map, proto
		Template string `json:"template"` // example: {"value":"$.booking_info.payments[0].type","tags":["payment"],"constraints":{"country_code":"$.country_code"}}
	}

	HTTPSinker struct {
		URL     string            `json:"url"`
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
	}

	FileSinker struct {
		FileName string `json:"file_name"`
	}
)
