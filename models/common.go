package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type MetaFields struct {
	Id          string     `bson:"_id"  json:"id"`
	CreatedTime *time.Time `bson:"createdTime,omitempty" json:"createdTime"`
}

func NewId() string {
	return bson.NewObjectId().Hex()
}

func NewMetaFields() MetaFields {
	now := time.Now()
	return MetaFields{
		Id: NewId(),
		CreatedTime: &now,
	}
}