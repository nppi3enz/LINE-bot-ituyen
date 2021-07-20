package db

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func Init(ctx context.Context) *firestore.Client {
	sa := option.WithCredentialsFile("google-credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	fmt.Println("run client : ")
	fmt.Println(client)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}
