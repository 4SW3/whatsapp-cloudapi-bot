package db

import "github.com/4d3v/gowtbot/pkg/models"

type DBRepo interface {
	SetBotControl(string, bool) error
	UpdateChat(string, models.DBMessage) error
	UpdateChatPlusBotControl(string, models.Message, bool) error
	CreateUser(models.Contact, models.Message) error
	GetUser(string) (string, bool, error)
	RemoveUser(string) error
	GetBlockedPhones() (map[string]bool, error)
	BlockOrUnblockUser(*models.AdminBlockUnblock) error
}
