package data

import (
	"sync"

	"github.com/4d3v/gowtbot/pkg/models"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var webhookPool = sync.Pool{
	New: func() any {
		return &models.WebhookBody{}
	},
}

var messagePool = sync.Pool{
	New: func() any {
		return &models.AdminMessage{}
	},
}

var blockUnblockPool = sync.Pool{
	New: func() any {
		return &models.AdminBlockUnblock{}
	},
}

func WebhookUnmarshalJSON(data []byte) (*models.WebhookBody, error) {
	webhookBody := webhookPool.Get().(*models.WebhookBody)
	if err := json.Unmarshal(data, &webhookBody); err != nil {
		return nil, err
	}
	return webhookBody, nil
}

func AdmMessageUnmarshalJSON(data []byte) (*models.AdminMessage, error) {
	msgBody := messagePool.Get().(*models.AdminMessage)
	if err := json.Unmarshal(data, &msgBody); err != nil {
		return nil, err
	}
	return msgBody, nil
}

func AdmBlockUnblockUnmarshalJSON(data []byte) (*models.AdminBlockUnblock, error) {
	blkunblkBody := blockUnblockPool.Get().(*models.AdminBlockUnblock)
	if err := json.Unmarshal(data, &blkunblkBody); err != nil {
		return nil, err
	}
	return blkunblkBody, nil
}
