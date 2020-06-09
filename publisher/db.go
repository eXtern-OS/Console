package publisher

import (
	"../db"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewDBCollection(collectionName string) (bool, *mongo.Collection) {
	return db.NewDatabaseCollection("dev", collectionName)
}

//User's id
func GetPublisherByUID(uid string) (bool, Publisher) {
	if t, c := NewDBCollection("publishers"); t {
		filter := bson.M{"maintainers_uids": uid}

		fmt.Println(filter)

		var res Publisher

		return c.FindOne(context.Background(), filter).Decode(&res) == nil, res
	} else {
		return false, Publisher{}
	}
}

func VerifyPublisherOwnsApp(appid, userid string) bool {
	if t, p := GetPublisherByUID(userid); t {
		for _, a := range p.Apps {
			if a == appid {
				return true
			}
		}
	}
	return false
}

func GetAppIds(uid string) (bool, []string) {
	if t, c := db.NewDBCollection("publishers"); t {
		filter := bson.M{
			"maintainers_uids": uid,
		}
		var res Publisher
		return c.FindOne(context.Background(), filter).Decode(&res) == nil, res.Apps
	}
	return false, []string{}
}
