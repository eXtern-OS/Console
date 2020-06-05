package app

import (
	"context"
	beatrix "github.com/eXtern-OS/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var URI = ""

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

	collection := client.Database("AppStore").Collection(collectionName)
	return true, collection
}

func GetPaidAppURL(id string) (int, string, string) {
	if t, c := NewDBCollection("apps", "GETPAIDAPPURL"); t {
		filter := bson.M{"app_id": id}
		var result Application

		err := c.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Println(err)
			log.Println("APPID", id)
			go beatrix.SendError("Error finding application by id", "GETPAIDAPPURL")
			return http.StatusInternalServerError, "", ""
		} else {
			slug := result.Slug
			filePath := result.Version.CurrentVersion.PackageURL

			if slug == "" || filePath == "" {
				go beatrix.SendError("ERROR: SLUG OR FILEPATH EMPTY", "GETPAIDAPPURL")
				return http.StatusInternalServerError, "", ""
			} else {
				return http.StatusOK, slug, filePath
			}
		}
	} else {
		go beatrix.SendError("ERROR CREATING NEW COLLECTION", "GETPAIDAPPURL")
		return http.StatusInternalServerError, "", ""
	}
}

func GetAppURL(id string) (bool, string) {
	if t, c := NewDBCollection("apps", "GETPAIDAPPURL"); t {
		filter := bson.M{"app_id": id}
		var result Application

		err := c.FindOne(context.Background(), filter).Decode(&result)
		if err != nil {
			log.Println(err)
			log.Println("APPID", id)
			go beatrix.SendError("Error finding application by id", "GETPAIDAPPURL")
			return false, ""
		} else {
			if !result.PaymentType.Free {
				//Sly ass trying to pass paid app url to get a free copy
				return false, ""
			} else {
				return true, result.Version.CurrentVersion.PackageURL
			}
		}
	} else {
		go beatrix.SendError("ERROR CREATING NEW COLLECTION", "GETPAIDAPPURL")
		return false, ""
	}
}

func GetAppByID(appid string) (bool, Application) {
	if t, c := NewDBCollection("apps", "GETAPPBYID"); t {
		filter := bson.M{"app_id": appid}
		var res Application
		err := c.FindOne(context.Background(), filter).Decode(&res)

		if err != nil {
			return false, Application{}
		}

		return true, res
	} else {
		return false, Application{}
	}
}
