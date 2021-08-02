package pubsub

type Msg struct {
	ID       string `json:"id,omitempty"`
	DID      string `json:"did"`
	RootHash string `json:"root_hash"`
}

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}

type Body struct {
	Message Msg `json:"Message"`
}
