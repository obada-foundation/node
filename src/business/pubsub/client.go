package pubsub

import (
	"context"
)

type Client interface {
	Publish(ctx context.Context, msg *Msg) (string, error)

	Subscribe(ctx context.Context) (*Msg, error)
}
