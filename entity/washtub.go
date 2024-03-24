package entity

import "time"

type PulseRequest struct {
	ChannelID string `json:"channel_id"`
	Address   string `json:"address"`
	Topic     string `json:"topic"`
	SinkType  string `json:"sink_type"`
	Status    string `json:"status"`
}

type PulseResponse struct {
	ResultStatus struct {
		Code    string      `json:"code"`
		Reason  string      `json:"reason"`
		Message interface{} `json:"message"`
	} `json:"result_status"`
	Data struct {
		ChannelID string    `json:"channel_id"`
		Address   string    `json:"address"`
		Topic     string    `json:"topic"`
		SinkType  string    `json:"sink_type"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}

type MessageRequest struct {
	ChannelID string `json:"channel_id"`
	Body      string `json:"body"`
	Status    string `json:"status"`
}

type MessageResponse struct {
	ResultStatus struct {
		Code    string      `json:"code"`
		Reason  string      `json:"reason"`
		Message interface{} `json:"message"`
	} `json:"result_status"`
	Data struct {
		ChannelID string    `json:"channel_id"`
		Body      string    `json:"body"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
}
