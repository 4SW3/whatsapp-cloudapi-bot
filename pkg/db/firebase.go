package db

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/4d3v/gowtbot/pkg/models"
)

const (
	fcol    = "whatsapp-messaging"
	fdoc    = "CoeK4Z4IeXk91bIbsWmS"
	bpdoc   = "gqKmfRXc2V8tLbaXY7uZ"
	timeout = 5 * time.Second
)

type firebaseRepo struct {
	FireApp *firebase.App
}

func NewFirebaseRepo(firebaseApp *firebase.App) DBRepo {
	return &firebaseRepo{FireApp: firebaseApp}
}

func (f *firebaseRepo) GetUser(phone string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return "", false, err
	}
	defer client.Close()

	dsnap, err := client.Collection(fcol).Doc(fdoc).Get(context.Background())
	if err != nil {
		return "", false, err
	}

	data := dsnap.Data()
	// check if the key which is the phone is in data
	if _, ok := data[phone]; !ok {
		return "", false, nil
	}

	botControl := data[phone].(map[string]any)["botControl"].(bool)

	return phone, botControl, nil
}

func (f *firebaseRepo) UpdateChat(phone string, message models.DBMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	now := time.Now().Unix()
	_, err = client.Collection(fcol).Doc(fdoc).Update(context.Background(), []firestore.Update{
		{
			Path: fmt.Sprintf("%s.%s", phone, "messages"),
			Value: firestore.ArrayUnion(map[string]any{
				"msgId":     message.GetID(),
				"body":      message.GetContent(),
				"sender":    message.GetSender(),
				"type":      message.GetType(),
				"createdAt": now,
				"updatedAt": now,
			}),
		},
		{
			Path:  fmt.Sprintf("%s.%s", phone, "updatedAt"),
			Value: now,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *firebaseRepo) UpdateChatPlusBotControl(phone string, message models.Message, controlVal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	now := time.Now().Unix()
	_, err = client.Collection(fcol).Doc(fdoc).Update(context.Background(), []firestore.Update{
		{
			Path: fmt.Sprintf("%s.%s", phone, "messages"),
			Value: firestore.ArrayUnion(map[string]any{
				"body":      fmt.Sprintf("Language set to %s", message.Interactive.ButtonReply.Title),
				"msgId":     message.ID,
				"sender":    message.From,
				"type":      message.Type,
				"createdAt": now,
				"updatedAt": now,
			}),
		},
		{
			Path:  fmt.Sprintf("%s.%s", phone, "botControl"),
			Value: controlVal,
		},
		{
			Path:  fmt.Sprintf("%s.%s", phone, "updatedAt"),
			Value: now,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *firebaseRepo) CreateUser(contact models.Contact, message models.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	now := time.Now().Unix()

	_, err = client.Collection(fcol).Doc(fdoc).Set(ctx, map[string]any{
		contact.WaID: map[string]any{
			"language":   "",
			"name":       contact.Profile.Name,
			"phone":      message.From,
			"botControl": true,
			"messages": []map[string]any{
				{
					"body":      message.Content,
					"msgId":     message.ID,
					"sender":    message.From,
					"type":      message.Type,
					"createdAt": now,
					"updatedAt": now,
				},
				{
					"body":      "choose language pt or en",
					"msgId":     fmt.Sprintf("%s-%d", "choose_language", now),
					"sender":    "bot",
					"type":      "interactive",
					"createdAt": now,
					"updatedAt": now,
				},
			},
			"createdAt": now,
			"updatedAt": now,
		},
	}, firestore.MergeAll)

	if err != nil {
		return err
	}

	return nil
}

func (f *firebaseRepo) RemoveUser(phone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// remove user from firebase using phone key
	_, err = client.Collection(fcol).Doc(fdoc).Update(context.Background(), []firestore.Update{
		{
			Path:  phone,
			Value: firestore.Delete,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *firebaseRepo) SetBotControl(phone string, controlVal bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Collection(fcol).Doc(fdoc).Update(context.Background(), []firestore.Update{
		{
			Path:  fmt.Sprintf("%s.%s", phone, "botControl"),
			Value: controlVal,
		},
		{
			Path:  fmt.Sprintf("%s.%s", phone, "updatedAt"),
			Value: time.Now().Unix(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (f *firebaseRepo) GetBlockedPhones() (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	dsnap, err := client.Collection(fcol).Doc(bpdoc).Get(context.Background())
	if err != nil {
		return nil, err
	}

	data := dsnap.Data()

	blockedPhones := make(map[string]bool)
	for k, v := range data {
		blockedPhones[k] = v.(bool)
	}

	return blockedPhones, nil
}

func (f *firebaseRepo) BlockOrUnblockUser(data *models.AdminBlockUnblock) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.FireApp.Firestore(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	fUpdate := firestore.Update{
		Path:  fmt.Sprintf("%s.%s", "blockedPhones", data.Phone),
		Value: true,
	}

	if !data.Block {
		fUpdate.Value = firestore.Delete
	}

	_, err = client.Collection(fcol).Doc(bpdoc).Update(context.Background(), []firestore.Update{fUpdate})
	if err != nil {
		return err
	}

	return nil
}
