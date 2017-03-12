package models

type User struct {
	MetaFields            `bson:",inline"`
	OpenId string         `bson:"openId"  json:"openId"`
}
