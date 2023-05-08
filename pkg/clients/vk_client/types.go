package vk_client

type LPServerResponse struct {
	Response LPServer
}

type LPServer struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type LongPoolResponse struct {
	Ts      string   `json:"ts"`
	Updates []Update `json:"updates"`
}

type Update struct {
	Type   string `json:"type"`
	Object Object `json:"object"`
}

type Object struct {
	Message Message
}

type Message struct {
	Date                  int          `json:"date"`
	FromId                int          `json:"from_id"`
	ID                    int          `json:"id"`
	Out                   int          `json:"out"`
	Attachments           []Attachment `json:"attachments"`
	ConversationMessageId int          `json:"conversation_message_id"`
	Keyboard              string       `json:"keyboard"`
	//"fwd_messages":[]
	//"important":false
	//"is_hidden":false
	PeerId   int    `json:"peer_id"`
	RandomId int    `json:"random_id"`
	Text     string `json:"text"`
}

type Attachment struct {
}

type MessageConfig struct {
	UserID   int    `json:"user_id"`
	RandomID int    `json:"random_id"`
	Message  string `json:"message"`
}
