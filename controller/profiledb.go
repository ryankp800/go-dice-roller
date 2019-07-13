package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)


func GetDBCollection() (*mongo.Collection, error) {
	// _ = os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	//mongodb://heroku_qkwm7vgb:tqug715ledj81r2gs24ajqj4kj@ds213255.mlab.com:13255/heroku_qkwm7vgb

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database("diceroller").Collection("users")
	return collection, nil
}
