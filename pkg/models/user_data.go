package models

// this is how the data is stored in firestore
// the data could be stored in the same way as it is received from the webhook
// but I choose to do this way so it makes it easier to deal with the data in the frontend
// because facebook webhook api has lots of nesting (see webhook_body.go)

type Msg struct {
	Body      string `json:"body"`
	MsgID     string `json:"msgId"`
	Sender    string `json:"sender"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type UserData struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	BotControl bool   `json:"botControl"`
	Language   string `json:"language"`
	Messages   []Msg  `json:"messages"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
}
