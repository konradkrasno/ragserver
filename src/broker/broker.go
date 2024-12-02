package broker

type Broker interface {
	Publish(string, []byte)
	Listen(string, func([]byte) error)
}
