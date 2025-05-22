package models

type WebhookRequest struct {
	Entry []Entry `json:"entry"`
}

type Entry struct {
    ID        string      `json:"id"`
    Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender  Sender  `json:"sender"`
	Message Message `json:"message"`
}

type Sender struct {
	ID string `json:"id"`
}

//=========facebook model=========

type FacebookReplyRequest struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}