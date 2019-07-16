package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client public mongo client
var Client *mongo.Client

// DiceCollection dice roller collection
var DiceCollection *mongo.Collection

// UserCollection dice roller collection
var UserCollection *mongo.Collection

// ConfigMongo sets up database
func ConfigMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// _ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	// mongodb://heroku_qkwm7vgb:tqug715ledj81r2gs24ajqj4kj@ds213255.mlab.com:13255/heroku_qkwm7vgb

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions); if err != nil {
		log.Fatal("this", err)
	}

	err = client.Ping(context.TODO(), nil); if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	DiceCollection = client.Database("heroku_qkwm7vgb").Collection("rolls")
}

func GetDiceRollByID(objectID string) DiceRoll {
	id, _ := primitive.ObjectIDFromHex(objectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tempDiceRoll DiceRoll
	if err := DiceCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tempDiceRoll); err != nil {
		log.Fatal(err)
	}

	return tempDiceRoll
}

func InsertDiceRoll(dieList DiceRoll) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := DiceCollection.InsertOne(ctx, dieList)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"new post created with id: %s\n",
		res.InsertedID.(primitive.ObjectID).Hex(),
	)

	GetDiceRollByID(res.InsertedID.(primitive.ObjectID).Hex())

}

func GetUsersCollection() (*mongo.Collection, error) {
	UserCollection = Client.Database("heroku_qkwm7vgb").Collection("users")
	return UserCollection, nil
}
