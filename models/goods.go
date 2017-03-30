package models

type Goods struct {
	MetaFields             `bson:",inline"`
	UserId string         `bson:"userId"  json:"userId"`
	Name   string         `bson:"name"  json:"name"`
	NamePY string         `bson:"name_py"  json:"name_py"`
	Status string         `bson:"status"  json:"status"`
}

const (
	GoodsStatusClose = "close"
)
