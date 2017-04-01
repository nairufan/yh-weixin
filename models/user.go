package models

type User struct {
	MetaFields            `bson:",inline"`
	OpenId   string       `bson:"openId"  json:"openId"`
	Nickname string       `bson:"nickname"  json:"nickname"`
	Sex      int          `bson:"sex"  json:"sex"` //1-man 2-女性
	City     string       `bson:"city"  json:"city"`
	Province string       `bson:"province"  json:"province"`
	Country  string       `bson:"country"  json:"country"`
	Avatar   string       `bson:"avatar"  json:"avatar"`
	UnionId  string       `bson:"unionid"  json:"unionid"`
}
