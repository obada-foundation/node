package queue

import (
	"context"
)

type MessageClient interface {
	Send(ctx context.Context, req *SendRequest) (string, error)
	// Receive Long polls given amount of messages from a queue.
	Receive(ctx context.Context, queueURL string) (*Message, error)
}
