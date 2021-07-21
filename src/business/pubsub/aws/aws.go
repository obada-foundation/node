package aws

import (
	"encoding/json"
	"github.com/obada-foundation/node/business/pubsub"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"time"
)

var _ pubsub.Client = Client{}

type Client struct {
	timeout time.Duration
	sqs *sqs.SQS
	sns *sns.SNS
	queueUrl string
	topicArn string
}

func NewClient(session *session.Session, timeout time.Duration, queueUrl, topicArn string) Client {
	return Client{
		timeout: timeout,
		sns: sns.New(session),
		sqs: sqs.New(session),
		queueUrl: queueUrl,
		topicArn: topicArn,
	}
}

func (c Client) Publish(ctx context.Context, msg *pubsub.Msg) (string, error) {
	//var r pubsub.SendRequest
	json, err := json.Marshal(msg)

	if err != nil {
		return "", err
	}

	jsonStr := string(json)

	messageGroupId := "OBADAUpdate"

	res, err := c.sns.PublishWithContext(ctx, &sns.PublishInput{
		Message:  &jsonStr,
		TopicArn: &c.topicArn,
		MessageGroupId: &messageGroupId,
	})

	if err != nil {
		return "", errors.Wrap(err, "cannot publish message to SNS")
	}

	return *res.MessageId, nil
}

func (c Client) Subscribe(ctx context.Context) (*pubsub.Msg, error) {
	var msg pubsub.Msg
	var b map[string]interface{}

	// timeout = WaitTimeSeconds + 5
	ctx, cancel := context.WithTimeout(ctx, time.Second*(20+5))
	defer cancel()

	res, err := c.sqs.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueUrl),
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
		msg.RootHash = fmt.Sprintf("%v", b["root_hash"])
	} else {
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", message)), &msg); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("pubsub :: subscribe :: unmarshal Message body %v: ", b))
		}
	}

	//msg.DID = b.Message.DID
	//msg.RootHash = b.Message.RootHash
	msg.ID = *res.Messages[0].MessageId

	return &msg, nil
}
