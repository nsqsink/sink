package contract

type Consumer interface {
	Run() error
	Stop() error
}
