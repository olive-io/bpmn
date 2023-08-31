package data

import (
	"context"

	"github.com/olive-io/bpmn/schema"
)

type runnerMessage interface {
	implementsRunnerMessage()
}

type getMessage struct {
	channel chan IItem
}

func (g getMessage) implementsRunnerMessage() {}

type putMessage struct {
	item    IItem
	channel chan struct{}
}

func (p putMessage) implementsRunnerMessage() {}

type Container struct {
	schema.ItemAwareInterface
	runnerChannel chan runnerMessage
	item          IItem
}

func NewContainer(ctx context.Context, itemAware schema.ItemAwareInterface) *Container {
	container := &Container{
		ItemAwareInterface: itemAware,
		runnerChannel:      make(chan runnerMessage, 1),
	}
	go container.run(ctx)
	return container
}

func (c *Container) Unavailable() bool {
	return false
}

func (c *Container) run(ctx context.Context) {
	for {
		select {
		case msg := <-c.runnerChannel:
			switch msg := msg.(type) {
			case getMessage:
				msg.channel <- c.item
			case putMessage:
				c.item = msg.item
				close(msg.channel)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Container) Get(ctx context.Context) <-chan IItem {
	ch := make(chan IItem)
	select {
	case c.runnerChannel <- getMessage{channel: ch}:
		return ch
	case <-ctx.Done():
		return nil
	}
}

func (c *Container) Put(ctx context.Context, item IItem) <-chan struct{} {
	ch := make(chan struct{})
	select {
	case c.runnerChannel <- putMessage{item: item, channel: ch}:
		return ch
	case <-ctx.Done():
		return nil
	}
}
