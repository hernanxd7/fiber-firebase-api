package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
)

func InitFirebase() {
	opt := option.WithCredentialsFile("config/credenciales.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error inicializando Firebase: %v\n", err)
		os.Exit(1)
	}

	FirebaseApp = app

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error inicializando Firestore: %v\n", err)
		os.Exit(1)
	}

	FirestoreClient = client
}
