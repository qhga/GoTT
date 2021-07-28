package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client
var colUsers *mongo.Collection
var colTypingTests *mongo.Collection
var colSurveys *mongo.Collection
var mongoCtx context.Context

func init() {
	// Mongo DB setup
	clientOptions := options.Client().ApplyURI("mongodb://" + config.DbUser + ":" + config.DbPass + "@" + config.DbHost + ":" + config.DbPort)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Cancel if Timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect to DB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	colUsers = client.Database("TT").Collection("users")
	colTypingTests = client.Database("TT").Collection("typingtests")
	colSurveys = client.Database("TT").Collection("surveys")
}
