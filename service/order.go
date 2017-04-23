package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/astaxie/beego"
)

const (
	collectionOrder = "order"
	collectionOrderItem = "orderItem"
)

func AddOrder(order *models.Order) *models.Order {
	order.MetaFields = models.NewMetaFields()
	order.Status = models.OrderStatusPending
	if order.UserId == "" {
		panic(apperror.NewInvalidParameterError("userId"))
	}
	if order.Name == "" {
		panic(apperror.NewInvalidParameterError("name"))
	}
	if order.Address == "" {
		panic(apperror.NewInvalidParameterError("address"))
	}
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionOrder, order)

	return order
}

func UpdateOrder(order *models.Order) *models.Order {
	if order.Id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}
	checkOrderStatus(order.Status)

	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionOrder, order.Id, order)
	return order
}

func GetOrderById(id string) *models.Order {
	session := mongo.Get()
	defer session.Close()
	order := &models.Order{}
	session.MustFindId(collectionOrder, id, order)
	return order
}

func AddOrderItem(orderItem *models.OrderItem) *models.OrderItem {
	orderItem.MetaFields = models.NewMetaFields()

	if orderItem.OrderId == "" {
		panic(apperror.NewInvalidParameterError("orderId"))
	}
	if orderItem.GoodsId == "" {
		panic(apperror.NewInvalidParameterError("goodsId"))
	}
	if orderItem.Quantity <= 0 {
		panic(apperror.NewInvalidParameterError("quantity"))
	}
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionOrderItem, orderItem)

	return orderItem
}

func AddOrderAgent(id string, agentId string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("orderId"))
	}
	order := GetOrderById(id)
	agents := order.Agents

	if agents == nil {
		agents = []*models.OrderAgent{}
	}

	find := false
	for _, ag := range agents {
		if ag.UpAgentId == agentId {
			find = true
			break
		}
	}

	if !find {
		a := &models.OrderAgent{
			UpAgentId: agentId,
			MetaFields: models.NewMetaFields(),
		}
		agents = append(agents, a)
		order.Agents = agents

		session := mongo.Get()
		defer session.Close()
		session.MustUpdateId(collectionOrder, order.Id, order)
	}
}

func RemoveOrderAgent(id string, agentId string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("orderId"))
	}
	order := GetOrderById(id)
	agents := order.Agents
	if agents != nil {
		index := -1
		for idx, a := range agents {
			if a.UpAgentId == agentId {
				index = idx
				break
			}
		}

		if index >= 0 {
			agents = append(agents[:index], agents[index + 1:]...)
			order.Agents = agents

			session := mongo.Get()
			defer session.Close()
			session.MustUpdateId(collectionOrder, order.Id, order)
		}
	}
}

func GetOrderItemById(id string) *models.OrderItem {
	session := mongo.Get()
	defer session.Close()
	orderItem := &models.OrderItem{}
	session.MustFindId(collectionOrderItem, id, orderItem)
	return orderItem
}

func RemoveOrderItemById(id string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}

	session := mongo.Get()
	defer session.Close()
	session.RemoveId(collectionOrderItem, id)
}

func RemoveAllOrderItems(id string) {
	if id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}

	session := mongo.Get()
	defer session.Close()
	session.RemoveAll(collectionOrderItem, bson.M{"orderId": id})
}

func UpdateOrderItem(orderItem *models.OrderItem) *models.OrderItem {
	if orderItem.Id == "" {
		panic(apperror.NewInvalidParameterError("id"))
	}
	item := GetOrderItemById(orderItem.Id)
	if item == nil {
		panic(apperror.NewInvalidParameterError("id"))
	}
	item.Quantity = orderItem.Quantity
	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionOrderItem, item.Id, item)
	return item
}

func GetOrders(userId string, offset int, limit int) []*models.Order {
	session := mongo.Get()
	defer session.Close()
	orders := []*models.Order{}

	option := mongo.Option{
		Sort: []string{"-createdTime"},
		Limit: &limit,
		Offset: &offset,
	}
	query := bson.M{}
	query["userId"] = userId
	query["status"] = bson.M{"$ne": models.OrderStatusClose}

	session.MustFindWithOptions(collectionOrder, query, option, &orders)
	return orders
}

/**
	agentId 上级代理商ID
 */
func GetAgentOrders(userId string, upAgentId string, offset int, limit int) []*models.Order {
	session := mongo.Get()
	defer session.Close()
	orders := []*models.Order{}

	option := mongo.Option{
		Sort: []string{"-createdTime"},
		Limit: &limit,
		Offset: &offset,
	}
	query := bson.M{}
	query["userId"] = userId
	query["upAgents"] = bson.M{"$elemMatch": bson.M{"upAgentId" : upAgentId}}

	session.MustFindWithOptions(collectionOrder, query, option, &orders)
	return orders
}

func GetOrdersByTimeRange(userId string, start time.Time, end time.Time) []*models.Order {
	session := mongo.Get()
	defer session.Close()
	orders := []*models.Order{}

	option := mongo.Option{
		Sort: []string{"-createdTime"},
	}
	query := bson.M{}
	query["userId"] = userId
	query["createdTime"] = bson.M{"$gte": start, "$lte": end}
	query["status"] = bson.M{"$ne": models.OrderStatusClose}

	beego.Info(query)
	session.MustFindWithOptions(collectionOrder, query, option, &orders)
	return orders
}

func GetOrderItems(ids []string) []*models.OrderItem {
	session := mongo.Get()
	defer session.Close()
	orderItems := []*models.OrderItem{}

	session.MustFind(collectionOrderItem, bson.M{"orderId": bson.M{"$in": ids}}, &orderItems)
	return orderItems
}

func GetOrderByTel(userId string, tel string, offset int, limit int) []*models.Order {
	if tel == "" {
		panic(apperror.NewInvalidParameterError("tel"))
	}
	session := mongo.Get()
	defer session.Close()
	orders := []*models.Order{}
	option := mongo.Option{
		Sort: []string{"-createdTime"},
		Limit: &limit,
		Offset: &offset,
	}
	query := bson.M{}
	query["userId"] = userId
	query["tel"] = tel
	query["status"] = bson.M{"$ne": models.OrderStatusClose}

	session.MustFindWithOptions(collectionOrder, query, option, &orders)
	return orders
}

func checkOrderStatus(status string) {
	if (status != "" && status != "pending" && status != "done" && status != "close") {
		panic(apperror.NewInvalidParameterError("status: pending, done, close"))
	}
}

func OrderStatistics(start time.Time, end time.Time) []*models.Statistic {
	results := []*models.Statistic{}
	statistics(start, end, collectionOrder, &results)
	return results
}

func OrderCount() int {
	session := mongo.Get()
	defer session.Close()
	return session.MustCount(collectionOrder)
}