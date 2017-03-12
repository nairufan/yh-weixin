package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionUser = "user"
)

func AddUser(user *models.User) *models.User {
	user.MetaFields = models.NewMetaFields()
	if user.OpenId == "" {
		panic(apperror.NewInvalidParameterError("openId"))
	}

	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionUser, user)

	return user
}

func GetUserByOpenId(openId string) *models.User {
	session := mongo.Get()
	defer session.Close()
	users := []*models.User{}
	session.MustFind(collectionUser, bson.M{"openId": openId}, &users)
	if len(users) > 0 {
		return users[0]
	}
	return nil
}