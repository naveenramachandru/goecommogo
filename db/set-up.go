package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)

	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println("Connected to mongose")

	return client

}

var Client *mongo.Client = DBSet()

func UserDeatils(client *mongo.Client, collectionName string) *mongo.Collection {
	var usercollection *mongo.Collection = client.Database("EcomProject").Collection(collectionName)

	return usercollection
}

func ProductDetails(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("EcomProject").Collection(collectionName)

	return collection
}

func PaymentDetails(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("EcomProject").Collection(collectionName)

	return collection
}

func OrderDeatils(client *mongo.Client, collectionName string) *mongo.Collection {
	var usercollection *mongo.Collection = client.Database("EcomProject").Collection(collectionName)

	return usercollection
}
