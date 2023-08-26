package models

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/4d3v/gowtbot/pkg/helpers"
)

// struct for message sent by admin
// example:
// {
//   "messaging_product": "whatsapp",
//   "preview_url": false,
//   "recipient_type": "individual",
//   "to": "5521972310861",
//   "type": "text",
//   "text": {
//     "body": "hello"
//   }
// }

// todo add contact, location
type AdminMessage struct {
	MessagingProduct string `json:"messaging_product"`
	PreviewURL       bool   `json:"preview_url"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Content          string `json:"content"`
	Text             struct {
		Body string `json:"body"`
	} `json:"text"`
	Image struct {
		Link string `json:"link"`
	} `json:"image"`
	Audio struct {
		Link string `json:"link"`
	} `json:"audio"`
	Video struct {
		Link string `json:"link"`
	} `json:"video"`
	Document struct {
		Link string `json:"link"`
	} `json:"document"`
	Sticker struct {
		Link string `json:"link"`
	} `json:"sticker"`
}

func (am AdminMessage) GetID() string {
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return fmt.Sprintf("%s_%s", am.To, time.Now().String())
	}
	return helpers.B2S(uuid)
}

func (am AdminMessage) GetContent() string {
	switch am.Type {
	case "text":
		return am.Text.Body
	case "image":
		return am.Image.Link
	case "audio":
		return am.Audio.Link
	case "video":
		return am.Video.Link
	case "document":
		return am.Document.Link
	case "sticker":
		return am.Sticker.Link
	default:
		return ""
	}
}

func (am AdminMessage) GetType() string {
	return am.Type
}

func (am AdminMessage) GetSender() string {
	return "admin"
}
