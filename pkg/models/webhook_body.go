package models

type (
	WebhookBody struct {
		Object string  `json:"object"`
		Entry  []Entry `json:"entry"`
	}

	Entry struct {
		ID      string   `json:"id"`
		Changes []Change `json:"changes"`
	}

	Change struct {
		Value Value  `json:"value"`
		Field string `json:"field"`
	}

	Value struct {
		MessagingProduct string    `json:"messaging_product"`
		Metadata         Metadata  `json:"metadata"`
		Contacts         []Contact `json:"contacts"`
		Messages         []Message `json:"messages"`
	}

	Metadata struct {
		DisplayPhoneNumber string `json:"display_phone_number"`
		PhoneNumberID      string `json:"phone_number_id"`
	}

	Contact struct {
		Profile Profile `json:"profile"`
		WaID    string  `json:"wa_id"`
	}

	Profile struct {
		Name string `json:"name"`
	}

	Message struct {
		From        string      `json:"from"`
		ID          string      `json:"id"`
		Timestamp   string      `json:"timestamp"`
		Type        string      `json:"type"`
		Content     string      `json:"content"`
		Text        Text        `json:"text"`
		Audio       Audio       `json:"audio"`
		Image       Image       `json:"image"`
		Video       Video       `json:"video"`
		Document    Document    `json:"document"`
		Interactive Interactive `json:"interactive"`
	}

	Interactive struct {
		Type        string      `json:"type"`
		ButtonReply ButtonReply `json:"button_reply"`
	}

	ButtonReply struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	Text struct {
		Body string `json:"body"`
	}

	Audio struct {
		ID string `json:"id"`
	}

	Image struct {
		ID string `json:"id"`
	}

	Video struct {
		ID string `json:"id"`
	}

	Document struct {
		ID string `json:"id"`
	}

	PayloadOptions struct {
		ClientPhone string
		LangID      string
	}
)

func (msg Message) GetID() string {
	return msg.ID
}

func (msg Message) GetContent() string {
	switch msg.Type {
	case "text":
		return msg.Text.Body
	case "image":
		return msg.Image.ID
	case "audio":
		return msg.Audio.ID
	case "video":
		return msg.Video.ID
	case "document":
		return msg.Document.ID
	default:
		return ""
	}
}

func (msg Message) GetType() string {
	return msg.Type
}

func (msg Message) GetSender() string {
	return msg.From
}
