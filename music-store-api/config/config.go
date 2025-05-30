package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://musicapp:2137@musicstore-cluster.yyqbx6d.mongodb.net/PAW-API-Database?retryWrites=true&w=majority&appName=musicstore-cluster"))
	if err != nil {
		log.Fatalf("Błąd łączenia z MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	DB = client.Database("PAW-API-Database")
	Client = client

	log.Println("MongoDB connected!")
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		log.Printf("Błąd rozłączania z MongoDB: %v", err)
	} else {
		log.Println("MongoDB disconnected!")
	}
}
