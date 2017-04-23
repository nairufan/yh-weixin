package models

type User struct {
	MetaFields            `bson:",inline"`
	OpenId   string       `bson:"openId"  json:"-"`
	Nickname string       `bson:"nickname"  json:"name"`
	Gender   int          `bson:"sex"  json:"sex"` //1-man 2-女性
	City     string       `bson:"city"  json:"city"`
	Province string       `bson:"province"  json:"province"`
	Country  string       `bson:"country"  json:"country"`
	Avatar   string       `bson:"avatar"  json:"avatar"`
	UnionId  string       `bson:"unionid"  json:"-"`
}

type UserAgent struct {
	MetaFields       `bson:",inline"`
	UserId  string   `bson:"userId"  json:"userId"`
	AgentId string   `bson:"agentId"  json:"agentId"`
	Name    string   `bson:"name"  json:"name"`
	Tel     string   `bson:"tel"  json:"tel"`
	Address string   `bson:"address"  json:"address"`
	Note    string   `bson:"note"  json:"note"`
	Avatar  string       `bson:"avatar"  json:"avatar"`
	Status  string   `bson:"status"  json:"status"`
}

type UserAgentBind struct {
	MetaFields       `bson:",inline"`
	UserId  string   `bson:"userId"  json:"userId"`
	AgentId string   `bson:"agentId"  json:"agentId"`
	Key     string   `bson:"key"  json:"key"`
}

const (
	UserAgentStatusDeleted = "deleted"
)