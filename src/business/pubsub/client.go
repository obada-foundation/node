package pubsub

import (
	"context"
)

// Client PubSub interface
type Client interface {
	Publish(ctx context.Context, msg *Msg) (string, error)

	Subscribe(ctx context.Context) (*Msg, error)
}
