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
	collectionUserAgent = "user_agent"
	collectionUserAgentBind = "user_agent_bind"
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

	if user.UnionId != "" {
		u.UnionId = user.UnionId
	}

	u.Nickname = user.Nickname
	u.Gender = user.Gender
	u.City = user.City
	u.Province = user.Province
	u.Country = user.Country
	u.Avatar = user.Avatar

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

func GetUserByIds(ids []string) []*models.User {
	session := mongo.Get()
	defer session.Close()
	users := []*models.User{}

	session.MustFind(collectionUser, bson.M{"_id": bson.M{"$in": ids}}, &users)
	return users
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

// agents
func AddUserAgent(id string, agentId string) *models.UserAgent {
	if id == "" {
		panic(apperror.NewInvalidParameterError("userId"))
	}
	if agentId == "" {
		panic(apperror.NewInvalidParameterError("agentId"))
	}

	agent := GetUserAgentByUserIdAndAgentId(id, agentId)
	if agent == nil {
		agentInfo := GetUserById(agentId)
		agent := &models.UserAgent{
			UserId: id,
			AgentId: agentId,
			Name: agentInfo.Nickname,
			Address: agentInfo.Country + agentInfo.Province + agentInfo.City,
			Avatar: agentInfo.Avatar,
		}
		agent.MetaFields = models.NewMetaFields()
		session := mongo.Get()
		defer session.Close()
		session.MustInsert(collectionUserAgent, agent)

		return agent
	}

	return agent
}

func UpdateUserAgent(agentId string, agent *models.UserAgent) {
	if agentId == "" {
		panic(apperror.NewInvalidParameterError("agentId"))
	}

	agentModel := GetUserAgentById(agentId)
	agentModel.Name = agent.Name
	agentModel.Address = agent.Address
	agentModel.Note = agent.Note
	agentModel.Tel = agent.Tel
	agentModel.Avatar = agent.Avatar
	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionUserAgent, agentModel.Id, agentModel)
}

func GetUserAgentByUserIdAndAgentId(userId string, agentId string) *models.UserAgent {
	agents := []*models.UserAgent{}
	session := mongo.Get()
	defer session.Close()
	session.MustFind(collectionUserAgent, bson.M{"userId": userId, "agentId": agentId}, &agents)

	if len(agents) > 0 {
		return agents[0]
	}

	return nil
}

func GetUserAgentById(id string) *models.UserAgent {
	session := mongo.Get()
	defer session.Close()
	agent := &models.UserAgent{}
	session.MustFindId(collectionUserAgent, id, agent)
	return agent
}

func GetUserAgentsByUserId(id string) []*models.UserAgent {
	session := mongo.Get()
	defer session.Close()
	agents := []*models.UserAgent{}

	query := bson.M{}
	query["userId"] = id
	query["status"] = bson.M{"$ne": models.UserAgentStatusDeleted}
	session.MustFind(collectionUserAgent, bson.M{"userId": id}, &agents)
	return agents
}

//user & agent bind service

func AddUserAgentBind(userId string, agentId string, key string) *models.UserAgentBind {
	if userId == "" {
		panic(apperror.NewInvalidParameterError("userId"))
	}
	if agentId == "" {
		panic(apperror.NewInvalidParameterError("agentId"))
	}
	if key == "" {
		panic(apperror.NewInvalidParameterError("key"))
	}

	session := mongo.Get()
	defer session.Close()

	bind := CheckUserAgentBind(userId, agentId)
	if bind == nil {
		bind = &models.UserAgentBind{
			UserId: userId,
			AgentId: agentId,
			Key: key,
		}
		bind.MetaFields = models.NewMetaFields()
		session.MustInsert(collectionUserAgentBind, bind)
	} else {
		bind.Key = key
		session.MustUpdateId(collectionUserAgentBind, bind.Id, bind)
	}
	return bind

}

func CheckUserAgentBind(userId string, agentId string) *models.UserAgentBind {
	binds := []*models.UserAgentBind{}
	session := mongo.Get()
	defer session.Close()
	session.MustFind(collectionUserAgentBind, bson.M{"userId": userId, "agentId": agentId}, &binds)

	if len(binds) > 0 {
		return binds[0]
	}

	return nil
}


func GetUserAgentBindByKey(key string) *models.UserAgentBind {
	binds := []*models.UserAgentBind{}
	session := mongo.Get()
	defer session.Close()
	session.MustFind(collectionUserAgentBind, bson.M{"key": key}, &binds)

	if len(binds) > 0 {
		return binds[0]
	}

	return nil
}