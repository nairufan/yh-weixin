package controllers

import (
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/apperror"
	"strconv"
	"reflect"
	"time"
	"github.com/nairufan/yh-weixin/utils"
	"strings"
)

var changeOrderFiledMap map[string]string

type OrderController struct {
	BaseController
}

type orderGoods struct {
	GoodsId  string     `json:"goodsId" validate:"required"`
	Quantity int        `json:"quantity" validate:"required"`
}

type createOrderRequest struct {
	CustomerName    string           `json:"customerName" validate:"required"`
	CustomerTel     string           `json:"customerTel" validate:"required"`
	CustomerAddress string           `json:"customerAddress"`
	GoodsList       []*orderGoods    `json:"goodsList" validate:"required"`
	TotalPrice      int              `json:"totalPrice"`
	Note            string           `json:"note"`
}

// @router /create [post]
func (o *OrderController) CreateOrder() {
	var request createOrderRequest
	o.Bind(&request)

	customer := mergeCustomer(request.CustomerName, request.CustomerTel, request.CustomerAddress, o.GetUserId())
	order := &models.Order{
		UserId: o.GetUserId(),
		Name: customer.Name,
		Tel: customer.Tel,
		Address: customer.Address,
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
	Id    string           `json:"id" validate:"required"`
	Field string           `json:"field" validate:"required"`
	Value string           `json:"value" validate:"required"`
}

// @router /update [post]
func (o *OrderController) UpdateOrder() {
	var request updateOrderRequest
	o.Bind(&request)
	filed := changeOrderFiledMap[request.Field]
	if filed == "" {
		panic(apperror.NewInvalidParameterError("valid field: name, tel, address, status, express, note, totalPrice"))
	}
	var value interface{}
	value = request.Value
	if filed == "TotalPrice" {
		value, _ = strconv.Atoi(request.Value)
	}
	order := updateOrderField(request.Id, filed, value)
	if filed == "Tel" {
		c := service.GetCustomerByTel(o.GetUserId(), request.Value)
		if c == nil {
			customer := &models.Customer{
				Name: order.Name,
				Tel: order.Tel,
				Address: order.Address,
				UserId: o.GetUserId(),
			}
			service.AddCustomer(customer)
		}
	}
	if filed == "Express" {
		order = updateOrderField(request.Id, "Status", models.OrderStatusDone)
	}
	o.Data["json"] = order
	o.ServeJSON()
}

type updateOrderItemsRequest struct {
	Id        string           `json:"orderId" validate:"required"`
	GoodsList []*orderGoods    `json:"goodsList" validate:"required"`
}

// @router /update-items [post]
func (o *OrderController) UpdateOrderItems() {
	var request updateOrderItemsRequest
	o.Bind(&request)

	service.RemoveAllOrderItems(request.Id)
	for _, goods := range request.GoodsList {
		orderItem := &models.OrderItem{
			OrderId: request.Id,
			GoodsId: goods.GoodsId,
			Quantity: goods.Quantity,
		}
		orderItem = service.AddOrderItem(orderItem)
	}
	o.Data["json"] = map[string]bool{"success": true}
	o.ServeJSON()
}

type deleteOrderItemRequest struct {
	Id string           `json:"id" validate:"required"`
}

// @router /remove-item [post]
func (o *OrderController) DeleteOrderItem() {
	var request deleteOrderItemRequest
	o.Bind(&request)

	service.RemoveOrderItemById(request.Id)
	o.Data["json"] = map[string]bool{"success": true}
	o.ServeJSON()
}

type updateOrderItemRequest struct {
	Id       string     `json:"id" `
	OrderId  string     `json:"orderId"`
	GoodsId  string     `json:"goodsId"`
	Quantity int        `json:"quantity" validate:"required"`
}

// @router /merge-item [post]
func (o *OrderController) UpdateOrderItem() {
	var request updateOrderItemRequest
	o.Bind(&request)

	item := &models.OrderItem{
		Quantity: request.Quantity,
	}
	if request.Id == "" {
		item.OrderId = request.OrderId
		item.GoodsId = request.GoodsId
		service.AddOrderItem(item)
	} else {
		item.Id = request.Id
		item = service.UpdateOrderItem(item)
	}

	o.Data["json"] = item
	o.ServeJSON()
}

func updateOrderField(id string, filedName string, fieldValue interface{}) *models.Order {
	order := service.GetOrderById(id)
	ps := reflect.ValueOf(order)
	val := ps.Elem()
	field := val.FieldByName(filedName)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(fieldValue))
	}
	newOrder := service.UpdateOrder(order)
	return newOrder
}

type orderListResponse struct {
	OrderList    []*models.Order                          `json:"orderList"`
	OrderItemMap map[string][]*models.OrderItem           `json:"orderItemMap"`
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

	for _, order := range orders {
		orderIds = append(orderIds, order.Id)
	}
	orderItems := service.GetOrderItems(orderIds)
	orderItemMap, goodsIds := ConvertOrderItemMap(orderItems)
	goodsList := service.GetGoodsByIds(goodsIds)
	response.OrderItemMap = orderItemMap
	response.GoodsMap = ConvertGoodsMap(goodsList)
	o.Data["json"] = response
	o.ServeJSON()
}

// @router /statistics [get]
func (o *OrderController) OrderStatistics() {
	response := models.StatisticResponse{}
	now := time.Now()
	start := now.AddDate(0, 0, -10)
	statistics := service.OrderStatistics(start, now)
	total := service.OrderCount()
	response.Statistics = statistics
	response.Total = total
	o.Data["json"] = response
	o.ServeJSON()
}

// @router /download [get]
func (o *OrderController) GetOrdersByRange() {
	start := o.GetString("start")
	end := o.GetString("end")
	if start == "" {
		start = time.Now().Format("2006-01-02")
	}
	if end == "" {
		end = time.Now().Format("2006-01-02")
	}
	s, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic(err)
	}
	e, err := time.Parse("2006-01-02", end)
	if err != nil {
		panic(err)
	}

	orders := service.GetOrdersByTimeRange(o.GetUserId(), s, e.AddDate(0, 0, 1))
	orderIds := []string{}

	for _, order := range orders {
		orderIds = append(orderIds, order.Id)
	}
	orderItems := service.GetOrderItems(orderIds)
	orderItemMap, goodsIds := ConvertOrderItemMap(orderItems)
	goodsList := service.GetGoodsByIds(goodsIds)
	goodsMap := ConvertGoodsMap(goodsList)
	fileName := ""
	if start == end {
		fileName = strings.Replace(start, "-", "", -1)
	} else {
		fileName = strings.Replace(start, "-", "", -1) + "～" + strings.Replace(end, "-", "", -1)
	}
	utils.Write(o.Ctx, getOrderRecords(orders, orderItemMap, goodsMap), fileName)
}

func mergeCustomer(customerName string, customerTel string, customerAddress string, userId string) *models.Customer {
	if customerName == "" {
		panic(apperror.NewParameterRequiredError("customerName"))
	}
	if customerTel == "" {
		panic(apperror.NewParameterRequiredError("customerTel"))
	}
	customer := service.GetCustomerByTel(userId, customerTel)
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

func ConvertGoodsMap(goodsList []*models.Goods) map[string]*models.Goods {
	goodsMap := map[string]*models.Goods{}
	for _, goods := range goodsList {
		goodsMap[goods.Id] = goods
	}

	return goodsMap
}

func getOrderRecords(orders []*models.Order, orderItemMap map[string][]*models.OrderItem, goodsMap map[string]*models.Goods) [][]string {
	results := [][]string{}
	results = append(results, []string{
		"电话号码",
		"姓名",
		"地址",
		"快递单号",
		"总价",
		"商品",
		"备注",
	})
	for _, order := range orders {
		record := []string{}
		record = append(record, order.Tel)
		record = append(record, order.Name)
		record = append(record, order.Address)
		record = append(record, order.Express)
		record = append(record, strconv.Itoa(order.TotalPrice))

		items := orderItemMap[order.Id]
		goods := ""
		for _, item := range items {
			goods += goodsMap[item.GoodsId].Name + "x" + strconv.Itoa(item.Quantity) + ", "
		}
		record = append(record, goods)
		record = append(record, order.Note)

		results = append(results, record)
	}

	return results
}

func init() {
	changeOrderFiledMap = map[string]string{
		"name": "Name",
		"tel": "Tel",
		"address": "Address",
		"status": "Status",
		"express": "Express",
		"note": "Note",
		"totalPrice": "TotalPrice",
	}
}