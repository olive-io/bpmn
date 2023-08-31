package event

type ISource interface {
	RegisterEventConsumer(IConsumer) error
}

type VoidSource struct{}

func (t VoidSource) RegisterEventConsumer(consumer IConsumer) (err error) {
	return
}
