package publisher

import (
	"context"

	"github.com/usesnipet/snipet/internal/runtime/event"
)

type Subscriber interface {
	Handle(context.Context, event.IEvent) error
}

type IPublisher interface {
	Publish(context.Context, event.IEvent) error
	Subscribe(s ...Subscriber)
}

type LocalPublisher struct {
	subscribers []Subscriber
}

func NewLocal() *LocalPublisher {
	return &LocalPublisher{
		subscribers: []Subscriber{},
	}
}

func (p *LocalPublisher) Subscribe(s ...Subscriber) {
	p.subscribers = append(p.subscribers, s...)
}

func (p *LocalPublisher) Publish(ctx context.Context, event event.IEvent) error {
	for _, s := range p.subscribers {
		if err := s.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
