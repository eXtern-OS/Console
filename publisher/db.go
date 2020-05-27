package publisher

import (
	"context"
	beatrix "github.com/eXtern-OS/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var URI string

func Init(mongouri string) {
	URI = mongouri
}

func NewDBCollection(collectionName, issuer string) (bool, *mongo.Collection) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error creating client for mongodb", issuer)
		return false, nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println(err)
		go beatrix.SendError("Error connecting with client to mongodb", issuer)
		return false, nil
	}

	collection := client.Database("dev").Collection(collectionName)
	return true, collection
}

//User's id
func GetPublisherByUID(uid string) (bool, Publisher) {
	if t, c := NewDBCollection("publishers", "GETPUBLISHER"); t {
		filter := bson.M{"maintainers_uids": uid}

		var res Publisher

		return c.FindOne(context.Background(), filter).Decode(res) == nil, res
	} else {
		return false, Publisher{}
	}
}
