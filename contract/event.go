package contract

// Event is an interface
// describing method for getting event detail
type Event interface {
	GetTopic() string
	GetSource() []string
}
