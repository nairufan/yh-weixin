package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
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
	if order.Tel == "" {
		panic(apperror.NewInvalidParameterError("tel"))
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