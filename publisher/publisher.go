package publisher

type Publisher struct {
	DisplayName     string   `bson:"display_name"     json:"display_name"`
	MaintainersUIDs []string `bson:"maintainers_uids" json:"maintainers_uids"`
	SumRatings      float64  `bson:"sum_ratings"      json:"sum_ratings"`
	Email           string   `bson:"email"            json:"email"`
	Address         string   `bson:"address"          json:"address"`
	Website         string   `bson:"website"          json:"website"`
	Verified        bool     `bson:"verified"         json:"verified"`
	UID             string   `bson:"uid"              json:"uid"`
}
