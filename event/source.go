package event

type Source interface {
	RegisterEventConsumer(Consumer) error
}

type VoidSource struct{}

func (t VoidSource) RegisterEventConsumer(consumer Consumer) (err error) {
	return
}
