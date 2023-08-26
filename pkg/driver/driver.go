package driver

import (
	"context"
	"path"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FirebaseInit() (*firebase.App, map[string]bool, error) {
	sa := option.WithCredentialsFile(path.Join("serviceAccount.json"))
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return nil, nil, err
	}

	blockedPhones, err := getBlockedPhones(app)
	if err != nil {
		return nil, nil, err
	}

	return app, blockedPhones, nil
}

func getBlockedPhones(f *firebase.App) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := f.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	fcol := "whatsapp-messaging"
	bpdoc := "gqKmfRXc2V8tLbaXY7uZ"

	dsnap, err := client.Collection(fcol).Doc(bpdoc).Get(context.Background())
	if err != nil {
		return nil, err
	}

	data := dsnap.Data()["blockedPhones"].(map[string]any)

	blockedPhones := make(map[string]bool)
	for k, v := range data {
		blockedPhones[k] = v.(bool)
	}

	blockedPhones["LOADED"] = true

	return blockedPhones, nil
}
