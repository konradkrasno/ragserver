package broker

type Broker interface {
	Publish(string, string, []byte)
	Listen(string, string, func([]byte) error)
}
