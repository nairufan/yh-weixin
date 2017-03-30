package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/nairufan/yh-weixin/utils"
)

const (
	collectionGoods = "goods"
)

func AddGoods(goods *models.Goods) *models.Goods {
	goods.MetaFields = models.NewMetaFields()
	if goods.Name == "" {
		panic(apperror.NewInvalidParameterError("name"))
	}
	if goods.UserId == "" {
		panic(apperror.NewInvalidParameterError("userId"))
	}
	goods.NamePY = utils.ConvertPY(goods.Name)
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionGoods, goods)

	return goods
}

func UpdateGoods(goods *models.Goods) *models.Goods {
	if goods.Id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}
	if goods.Name == "" {
		panic(apperror.NewInvalidParameterError("name"))
	}
	g := GetGoodsById(goods.Id)
	g.Name = goods.Name
	g.NamePY = utils.ConvertPY(goods.Name)
	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionGoods, goods.Id, g)
	return g
}

func GetGoodsById(id string) *models.Goods {
	session := mongo.Get()
	defer session.Close()
	goods := &models.Goods{}
	session.MustFindId(collectionGoods, id, goods)
	return goods
}

func RemoveGoodsById(id string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}

	goods := GetGoodsById(id)
	if goods == nil {
		panic(apperror.NewResourceNotFoundError("goods"))
	}
	goods.Status = models.GoodsStatusClose

	session := mongo.Get()
	defer session.Close()

	session.MustUpdateId(collectionGoods, id, goods)
}

func GetGoods(userId string, offset int, limit int) []*models.Goods {
	session := mongo.Get()
	defer session.Close()
	goods := []*models.Goods{}

	option := mongo.Option{
		Sort: []string{"+name_py"},
		Limit: &limit,
		Offset: &offset,
	}
	query := bson.M{}
	query["userId"] = userId
	query["status"] = bson.M{"$ne": models.GoodsStatusClose}

	session.MustFindWithOptions(collectionGoods, query, option, &goods)
	return goods
}

func GetGoodsByIds(ids []string) []*models.Goods {
	session := mongo.Get()
	defer session.Close()
	goods := []*models.Goods{}

	session.MustFind(collectionGoods, bson.M{"_id": bson.M{"$in": ids}}, &goods)
	return goods
}

func GoodsStatistics(start time.Time, end time.Time) []*models.Statistic {
	results := []*models.Statistic{}
	statistics(start, end, collectionGoods, &results)
	return results
}

func GoodsCount() int {
	session := mongo.Get()
	defer session.Close()
	return session.MustCount(collectionGoods)
}

func GetAllGoods() []*models.Goods {
	session := mongo.Get()
	defer session.Close()
	goods := []*models.Goods{}
	session.MustFind(collectionGoods, bson.M{}, &goods)
	return goods
}