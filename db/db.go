package db

import (
	"context"
	"github.com/masci/flickr"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

var URI = ""

var Client MongoClient

func Init(mongouri, fapi, fsecret string) {
	URI = mongouri
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	Client.Mutex.Lock()
	Client.Client = client
	Client.Mutex.Unlock()

	fclient := flickr.NewFlickrClient(fapi, fsecret)
	FClient.Mutex.Lock()
	FClient.Client = fclient
	FClient.Mutex.Unlock()
}

type MongoClient struct {
	Mutex  sync.Mutex
	Client *mongo.Client
}

func (c *MongoClient) Reload() {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	Client.Mutex.Lock()
	Client.Client = client
	Client.Mutex.Unlock()
}

func NewDBCollection(collectionName string) (bool, *mongo.Collection) {
	Client.Mutex.Lock()
	collection := Client.Client.Database("dev").Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}

func NewDatabaseCollection(database, collectionName string) (bool, *mongo.Collection) {
	Client.Mutex.Lock()
	collection := Client.Client.Database(database).Collection(collectionName)
	Client.Mutex.Unlock()
	return true, collection
}
