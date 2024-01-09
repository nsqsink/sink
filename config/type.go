package config

type (
	App struct {
		LogLevel  string     `json:"log_level"`
		Consumers []Consumer `json:"consumers"`
	}

	Consumer struct {
		ID          string `json:"id"`     // will be channel name
		Topic       string `json:"topic"`  // source topic
		Source      string `json:"source"` // list source lookupd, separated by comma for multiple value
		MaxAttempt  int    `json:"max_attempt"`
		MaxInFlight int    `json:"max_in_flight"`
		Concurrent  int    `json:"concurrent"`
		Sinker      Sinker `json:"sinker"`
		Active      bool   `json:"active"`
	}

	Sinker struct {
		Type   string     `json:"type"` //sinker type
		Parser Parser     `json:"parser"`
		HTTP   HTTPSinker `json:"http"`
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
)
