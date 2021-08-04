package pubsub

// Msg message
type Msg struct {
	ID       string `json:"id,omitempty"`
	DID      string `json:"did"`
	Checksum string `json:"checksum"`
}

// SendRequest request to SNS topic
type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

// Attribute SNS topic attribute
type Attribute struct {
	Key   string
	Value string
	Type  string
}

// Body SQS body
type Body struct {
	Message Msg `json:"Message"`
}
