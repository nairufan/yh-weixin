package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
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

func UpdateUser(user *models.User) *models.User {
	if user.Id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}
	u := GetUserById(user.Id)
	u.UnionId = user.UnionId

	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionUser, u.Id, u)
	return u
}

func GetUserById(id string) *models.User {
	session := mongo.Get()
	defer session.Close()
	user := &models.User{}
	session.MustFindId(collectionUser, id, user)
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

func GetUserByUnionId(unionId string) *models.User {
	session := mongo.Get()
	defer session.Close()
	users := []*models.User{}
	session.MustFind(collectionUser, bson.M{"unionid": unionId}, &users)
	if len(users) > 0 {
		return users[0]
	}
	return nil
}

func UserStatistics(start time.Time, end time.Time) []*models.Statistic {
	results := []*models.Statistic{}
	statistics(start, end, collectionUser, &results)
	return results
}

func UserCount() int {
	session := mongo.Get()
	defer session.Close()
	return session.MustCount(collectionUser)
}