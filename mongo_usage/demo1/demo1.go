package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
	"time"
)

func main() {
	var (
		client        *mongo.Client
		err           error
		database      *mongo.Database
		collection    *mongo.Collection
		clientOptions *options.ClientOptions
	)
	clientOptions = options.Client()
	clientOptions.SetConnectTimeout(time.Duration(5 * time.Second))
	if client, err = mongo.Connect(context.TODO(), "mongodb://47.75.179.127:27017"); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("mydb")

	collection = database.Collection("my_collection")

	collection = collection
}
