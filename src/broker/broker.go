package broker

type Broker interface {
	Publish(string, string, []byte) error
	Listen(string, string, func([]byte) error) error
}
