package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var (
	Item *mongo.Collection
	Seller  *mongo.Collection
	User    *mongo.Collection
)

var Mutex = sync.Mutex{}

func Init() (*mongo.Database, error) {
	fmt.Println("Running the database.Init function")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/dc-ecommrce")
	client, err := mongo.Connect(context.Background(), clientOptions)
	
	if err != nil {
		return nil, err
	}

	Item = client.Database("dc-ecommerce").Collection("items")
	Seller = client.Database("dc-ecommerce").Collection("sellers")
	User = client.Database("dc-ecommerce").Collection("users")

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("dc-ecommerce"), nil
}
