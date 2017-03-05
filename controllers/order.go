package controllers

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/apperror"
	"strconv"
)

type OrderController struct {
	BaseController
}

type orderGoods struct {
	GoodsId  string     `json:"goodsId" validate:"required"`
	Quantity int        `json:"quantity" validate:"required"`
}

type createOrderRequest struct {
	CustomerId      string           `json:"customerId"`
	CustomerName    string           `json:"customerName"`
	CustomerTel     string           `json:"customerTel"`
	CustomerAddress string           `json:"customerAddress"`
	GoodsList       []*orderGoods    `json:"goodsList" validate:"required"`
	TotalPrice      int              `json:"totalPrice"`
	Note            string           `json:"note"`
}

// @router /create [post]
func (o *OrderController) CreateOrder() {
	var request createOrderRequest
	o.Bind(&request)

	customer := mergeCustomer(request.CustomerId, request.CustomerName, request.CustomerTel, request.CustomerAddress, o.GetUserId())
	order := &models.Order{
		UserId: o.GetUserId(),
		CustomerId: customer.Id,
		Tel: customer.Tel,
		TotalPrice: request.TotalPrice,
		Note: request.Note,
	}
	order = service.AddOrder(order)
	for _, goods := range request.GoodsList {
		orderItem := &models.OrderItem{
			OrderId: order.Id,
			GoodsId: goods.GoodsId,
			Quantity: goods.Quantity,
		}
		orderItem = service.AddOrderItem(orderItem)
	}
	o.Data["json"] = map[string]bool{"success": true}
	o.ServeJSON()
}

type updateOrderRequest struct {
	Id      string           `json:"id" validate:"required"`
	Status  string           `json:"status" validate:"required"`
	Express string           `json:"express" validate:"required"`
}

// @router /update [post]
func (o *OrderController) UpdateOrder() {
	var request updateOrderRequest
	o.Bind(&request)

	order := &models.Order{
		Status: request.Status,
		Express: request.Express,

	}
	order.Id = request.Id
	order = service.UpdateOrder(order)
	o.Data["json"] = order
	o.ServeJSON()
}

type orderListResponse struct {
	OrderList    []*models.Order                          `json:"orderList"`
	OrderItemMap map[string][]*models.OrderItem           `json:"orderItemMap"`
	CustomerMap  map[string]*models.Customer                   `json:"customerMap"`
	GoodsMap     map[string]*models.Goods                 `json:"goodsMap"`
}

// @router /list [get]
func (o *OrderController) GetOrders() {
	offset := o.GetString("offset")
	limit := o.GetString("limit")
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)
	response := &orderListResponse{}

	orders := service.GetOrders(o.GetUserId(), offsetInt, limitInt)
	response.OrderList = orders
	orderIds := []string{}
	customerIds := []string{}
	for _, order := range orders {
		orderIds = append(orderIds, order.Id)
		customerIds = append(customerIds, order.CustomerId)
	}
	orderItems := service.GetOrderItems(orderIds)
	customers := service.GetCustomerByIds(customerIds)
	orderItemMap, goodsIds := ConvertOrderItemMap(orderItems)
	goodsList := service.GetGoodsByIds(goodsIds)
	response.OrderItemMap = orderItemMap
	response.CustomerMap = ConvertCustomerMap(customers)
	response.GoodsMap = ConvertGoodsMap(goodsList)
	o.Data["json"] = response
	o.ServeJSON()
}

func mergeCustomer(customerId string, customerName string, customerTel string, customerAddress string, userId string) *models.Customer {
	if customerId != "" {
		c := service.GetCustomerById(customerId)
		if c == nil {
			panic(apperror.NewInvalidParameterError("customerId"))
		}
		return c
	}
	if customerName == "" {
		panic(apperror.NewParameterRequiredError("customerName"))
	}
	if customerTel == "" {
		panic(apperror.NewParameterRequiredError("customerTel"))
	}
	if customerAddress == "" {
		panic(apperror.NewParameterRequiredError("customerAddress"))
	}
	customer := service.GetCustomerByTel(customerTel)
	if customer != nil {
		customer.Name = customerName
		customer.Tel = customerTel
		customer.Address = customerAddress

		return service.UpdateCustomer(customer)
	}

	customer = &models.Customer{
		Name: customerName,
		Tel: customerTel,
		Address: customerAddress,
		UserId: userId,
	}
	return service.AddCustomer(customer)
}

func ConvertOrderItemMap(orderItems []*models.OrderItem) (map[string][]*models.OrderItem, []string) {
	orderItemMap := map[string][]*models.OrderItem{}
	goodsIds := []string{}
	for _, orderItem := range orderItems {
		items := orderItemMap[orderItem.OrderId]
		if items == nil {
			items = []*models.OrderItem{}
		}
		items = append(items, orderItem)
		orderItemMap[orderItem.OrderId] = items
		goodsIds = append(goodsIds, orderItem.GoodsId)
	}

	return orderItemMap, goodsIds
}

func ConvertCustomerMap(customers []*models.Customer) map[string]*models.Customer {
	customerMap := map[string]*models.Customer{}
	for _, customer := range customers {
		customerMap[customer.Id] = customer
	}

	return customerMap
}

func ConvertGoodsMap(goodsList []*models.Goods) map[string]*models.Goods {
	goodsMap := map[string]*models.Goods{}
	for _, goods := range goodsList {
		goodsMap[goods.Id] = goods
	}

	return goodsMap
}