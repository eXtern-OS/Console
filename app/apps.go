package app

import (
	"context"
	beatrix "github.com/eXtern-OS/Beatrix"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func GetAppsByCategories(categories []string, limit int, collection *mongo.Collection) ExportedApplications {
	var result ExportedApplications

	for _, c := range categories {
		filter := bson.M{"category": c}

		//TODO: Not to reindex all apps better create something like genesis block
		// Hmmmmm
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			log.Println(err)
			go beatrix.SendError("Error finding queries in db", "GETAPPSBYCATEGORIES")
			return ExportedApplications{}
		}
		var res []Application
		err = cursor.Decode(&res)
		if err != nil {
			log.Println(err)
			go beatrix.SendError("Error decoding value", "GETAPPSBYCATEGORIES")
			return ExportedApplications{}
		}
		result.Apps = append(result.Apps, ExportSlice(res[:limit])...)
	}

	if len(result.Apps) == 0 {
		log.Println("ERROR: STILL 0 EXPORTEDAPPLICATIONS")
		go beatrix.SendError("ERROR: STILL 0 ELEMENTS", "GETAPPSBYCATEGORIES")
		return ExportedApplications{}
	}
	return result
}
