package models

type DBMessage interface {
	GetID() string
	GetContent() string
	GetType() string
	GetSender() string
}
