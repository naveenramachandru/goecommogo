package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QEJM2c9l5RNOxdEh
// 122.172.87.33/32
// mongodb+srv://naveenramachandru77:QEJM2c9l5RNOxdEh@cluster0.7wun6.mongodb.net/

func DBSet() *mongo.Client {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://naveenramachandru77:QEJM2c9l5RNOxdEh@cluster0.7wun6.mongodb.net/"))

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
