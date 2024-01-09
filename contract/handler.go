package contract

type Handler interface {
	Handle(msg Messager) error
}
