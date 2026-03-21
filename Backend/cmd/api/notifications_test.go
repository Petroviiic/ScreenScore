package main

import (
	"context"
	"testing"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func TestFirebaseInit(t *testing.T) {
	opt := option.WithCredentialsFile("../../screenscore_firebase_key.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		t.Fatalf("Firebase init failed: %v", err)
	}

	_, err = app.Messaging(context.Background())
	if err != nil {
		t.Fatalf("Messaging client failed: %v", err)
	}
}
