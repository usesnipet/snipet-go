package runtime

import "context"

type Subscriber interface {
	Handle(context.Context, IEvent) error
}

type IPublisher interface {
	Publish(context.Context, IEvent) error
	Subscribe(s ...Subscriber)
}

type LocalPublisher struct {
	subscribers []Subscriber
}

func NewLocalPublisher() *LocalPublisher {
	return &LocalPublisher{
		subscribers: []Subscriber{},
	}
}

func (p *LocalPublisher) Subscribe(s ...Subscriber) {
	p.subscribers = append(p.subscribers, s...)
}

func (p *LocalPublisher) Publish(ctx context.Context, event IEvent) error {
	for _, s := range p.subscribers {
		if err := s.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
