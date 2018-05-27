package core

import "gopkg.in/mgo.v2/bson"

// MgoDBQuery ...
type MgoDBQuery struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func CountAnyCollectionWithQuery(collectionName string, query bson.M) (int, error) {
	db := MgoDb{}
	db.InitByRevelConfig()
	defer db.Close()
	return db.C(collectionName).Find(query).Count()
}
