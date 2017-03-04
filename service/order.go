package service

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/apperror"
	"github.com/nairufan/yh-weixin/db/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionOrder = "order"
	collectionOrderItem = "orderItem"
)

func AddOrder(order *models.Order) *models.Order {
	order.MetaFields = models.NewMetaFields()
	order.Status = models.OrderStatusPending
	if order.CustomerId == "" {
		apperror.NewInvalidParameterError("customerId")
	}
	if order.UserId == "" {
		apperror.NewInvalidParameterError("userId")
	}
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionOrder, order)

	return order
}

func UpdateOrder(order *models.Order) *models.Order {
	if order.Id == "" {
		apperror.NewInvalidParameterError("id")
	}
	checkOrderStatus(order.Status)
	o := GetOrderById(order.Id)
	o.Express = order.Express
	o.Status = order.Status

	session := mongo.Get()
	defer session.Close()
	session.MustUpdateId(collectionOrder, o.Id, o)
	return o
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
		apperror.NewInvalidParameterError("orderId")
	}
	if orderItem.GoodsId == "" {
		apperror.NewInvalidParameterError("goodsId")
	}
	if orderItem.Quantity <= 0 {
		apperror.NewInvalidParameterError("quantity")
	}
	session := mongo.Get()
	defer session.Close()
	session.MustInsert(collectionOrderItem, orderItem)

	return orderItem
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
	session.MustFindWithOptions(collectionOrder, bson.M{"userId": userId}, option, &orders)
	return orders
}

func GetOrderItems(ids []string) []*models.OrderItem {
	session := mongo.Get()
	defer session.Close()
	orderItems := []*models.OrderItem{}

	session.MustFind(collectionOrderItem, bson.M{"orderId": bson.M{"$in": ids}}, &orderItems)
	return orderItems
}

func checkOrderStatus(status string) {
	if (status != "" && status != "pending" && status != "done" && status != "close") {
		panic(apperror.NewInvalidParameterError("status: pending, done, close"))
	}
}