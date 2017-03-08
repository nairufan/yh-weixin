package models

type Order struct {
	MetaFields                `bson:",inline"`
	UserId     string         `bson:"userId"  json:"userId"`
	Name       string         `bson:"name"  json:"name"`
	Tel        string         `bson:"tel"  json:"tel"`
	Address    string         `bson:"address"  json:"address"`
	Status     string         `bson:"status"  json:"status"`
	Express    string         `bson:"express"  json:"express"`
	Note       string         `bson:"note"  json:"note"`
	TotalPrice int            `bson:"totalPrice"  json:"totalPrice"`
}

type OrderItem struct {
	MetaFields                 `bson:",inline"`
	OrderId  string            `bson:"orderId"  json:"orderId"`
	GoodsId  string            `bson:"goodsId"  json:"goodsId"`
	Quantity int               `bson:"quantity"  json:"quantity"`
}

const (
	OrderStatusPending = "pending"
	OrderStatusDone = "done"
	OrderStatusClose = "close"
)
