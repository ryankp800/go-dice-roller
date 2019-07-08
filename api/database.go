package api

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

//Client public mongo client
var Client *mongo.Client

//Collection dice roller collection
var Collection *mongo.Collection

//ConfigMongo sets up database
func ConfigMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	//mongodb://heroku_qkwm7vgb:tqug715ledj81r2gs24ajqj4kj@ds213255.mlab.com:13255/heroku_qkwm7vgb
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("this", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	Collection = client.Database("diceroller").Collection("rolls")
}

func getDiceRollByID(objectID string) DiceRoll {
	id, _ := primitive.ObjectIDFromHex(objectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tempDiceRoll DiceRoll
	if err := Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tempDiceRoll); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("post : %+v\n", tempDiceRoll)

	return tempDiceRoll
}

func insertDiceRoll(dieList DiceRoll) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := Collection.InsertOne(ctx, dieList)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"new post created with id: %s\n",
		res.InsertedID.(primitive.ObjectID).Hex(),
	)

	getDiceRollByID(res.InsertedID.(primitive.ObjectID).Hex())

}
