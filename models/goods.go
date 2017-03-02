package models

type Goods struct {
	MetaFields             `bson:",inline"`
	UserId string         `bson:"userId"  json:"userId"`
	Name   string         `bson:"name"  json:"name"`
}

