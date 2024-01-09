package contract

type Parser interface {
	// parse given data into given format on the module
	Parse(data []byte) (result []byte, err error)
}
