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

var database = "heroku_qkwm7vgb"

// ConfigMongo sets up database
func ConfigMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := os.Getenv("MONGODB_URI")

	if url == "" {
		url  =  "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, clientOptions); if err != nil {
		log.Fatal("this", err)
	}

	err = client.Ping(context.TODO(), nil); if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	DiceCollection = client.Database(database).Collection("rolls")
	UserCollection = client.Database(database).Collection("users")
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

func insertDiceRoll(dieList DiceRoll) {
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
