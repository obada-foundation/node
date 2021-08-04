package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/obada-foundation/node/business/pubsub"
	"github.com/pkg/errors"
	"time"
)

var _ pubsub.Client = Client{}

// Client AWS PubSub client
type Client struct {
	timeout  time.Duration
	sqs      *sqs.SQS
	sns      *sns.SNS
	queueURL string
	topicArn string
}

// NewClient creates AWS PubSub client
func NewClient(sess *session.Session, timeout time.Duration, queueURL, topicArn string) Client {
	return Client{
		timeout:  timeout,
		sns:      sns.New(sess),
		sqs:      sqs.New(sess),
		queueURL: queueURL,
		topicArn: topicArn,
	}
}

// Publish changes to SNS topic
func (c Client) Publish(ctx context.Context, msg *pubsub.Msg) (string, error) {
	bytes, err := json.Marshal(msg)

	if err != nil {
		return "", err
	}

	jsonStr := string(bytes)

	messageGroupID := "OBADAUpdate"

	res, err := c.sns.PublishWithContext(ctx, &sns.PublishInput{
		Message:        &jsonStr,
		TopicArn:       &c.topicArn,
		MessageGroupId: &messageGroupID,
	})

	if err != nil {
		return "", errors.Wrap(err, "cannot publish message to SNS")
	}

	return *res.MessageId, nil
}

// Subscribe receives changes from SQS
func (c Client) Subscribe(ctx context.Context) (*pubsub.Msg, error) {
	var msg pubsub.Msg
	var b map[string]interface{}

	ctx, cancel := context.WithTimeout(ctx, time.Second*(20+5))
	defer cancel()

	res, err := c.sqs.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(20),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if len(res.Messages) == 0 {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(*res.Messages[0].Body), &b); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("pubsub :: subscribe :: unmarshal SQS body %v: ", *res.Messages[0].Body))
	}

	message, ok := b["Message"]

	if !ok {
		msg.DID = fmt.Sprintf("%v", b["did"])
		msg.Checksum = fmt.Sprintf("%v", b["checksum"])
	} else if err := json.Unmarshal([]byte(fmt.Sprintf("%v", message)), &msg); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("pubsub :: subscribe :: unmarshal Message body %v: ", b))
	}

	msg.ID = *res.Messages[0].MessageId

	return &msg, nil
}
