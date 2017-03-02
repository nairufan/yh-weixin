package models

type Customer struct {
	MetaFields             `bson:",inline"`
	UserId  string         `bson:"userId"  json:"userId"`
	Name    string         `bson:"name"  json:"name"`
	Tel     string         `bson:"tel"  json:"tel"`
	Address string         `bson:"address"  json:"address"`
}
