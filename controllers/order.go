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
	CustomerTel     string           `json:"customerTel"`
	CustomerAddress string           `json:"customerAddress"`
	GoodsList       []*orderGoods    `json:"goodsList"`
	TotalPrice      int              `json:"totalPrice"`
	Note            string           `json:"note"`
}

// @router /create [post]
func (o *OrderController) CreateOrder() {
	var request createOrderRequest
	o.Bind(&request)

	customer := &models.Customer{}
	if request.CustomerTel != "" {
		customer = mergeCustomer(request.CustomerName, request.CustomerTel, request.CustomerAddress, o.GetUserId())
	} else {
		customer.Name = request.CustomerName
		customer.Address = request.CustomerAddress
	}

	order := &models.Order{
		UserId: o.GetUserId(),
		Name: customer.Name,
		Tel: customer.Tel,
		Address: customer.Address,
		TotalPrice: request.TotalPrice,
		Note: request.Note,
	}
	order = service.AddOrder(order)
	if request.GoodsList != nil {
		for _, goods := range request.GoodsList {
			orderItem := &models.OrderItem{
				OrderId: order.Id,
				GoodsId: goods.GoodsId,
				Quantity: goods.Quantity,
			}
			service.AddOrderItem(orderItem)
		}
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
	order := updateOrderField(o.GetUserId(), request.Id, filed, value)
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
		order = updateOrderField(o.GetUserId(), request.Id, "Status", models.OrderStatusDone)
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

type agentOrderItemRequest struct {
	AgentId         string            `json:"upAgentId" validate:"required"`
	AddedOrderIds   []string          `json:"addedOrderIds"`
	RemovedOrderIds []string          `json:"removedOrderIds"`
}

// @router /update-agents [post]
func (o *OrderController) UpdateOrderAgents() {
	var request agentOrderItemRequest
	o.Bind(&request)

	if request.AddedOrderIds != nil {
		for _, id := range request.AddedOrderIds {
			service.AddOrderAgent(o.GetUserId(), id, request.AgentId)
		}
	}

	if request.RemovedOrderIds != nil {
		for _, id := range request.RemovedOrderIds {
			service.RemoveOrderAgent(o.GetUserId(), id, request.AgentId)
		}
	}

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

func updateOrderField(userId string, id string, filedName string, fieldValue interface{}) *models.Order {
	order := service.GetOrderById(id)
	ps := reflect.ValueOf(order)
	val := ps.Elem()
	field := val.FieldByName(filedName)
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(fieldValue))
	}
	newOrder := service.UpdateOrder(userId, order)
	return newOrder
}

type OrderListResponse struct {
	OrderList    []*models.Order                          `json:"orderList"`
	OrderItemMap map[string][]*models.OrderItem           `json:"orderItemMap"`
	GoodsMap     map[string]*models.Goods                 `json:"goodsMap"`
}

// @router /list [get]
func (o *OrderController) GetOrders() {
	offset := o.GetString("offset")
	limit := o.GetString("limit")
	o.Data["json"] = orderList(o.GetUserId(), offset, limit)
	o.ServeJSON()
}

func orderList(userId string, offset string, limit string) *OrderListResponse {
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)
	response := &OrderListResponse{}

	orders := service.GetOrders(userId, offsetInt, limitInt)
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

	return response
}

type selectOrder struct {
	*models.Order
	IsSelected bool     `json:"isSelected"`
}

type selectOrderListResponse struct {
	OrderList    []*selectOrder                           `json:"orderList"`
	OrderItemMap map[string][]*models.OrderItem           `json:"orderItemMap"`
	GoodsMap     map[string]*models.Goods                 `json:"goodsMap"`
}

// @router /list-select [get]
func (o *OrderController) GetSelectOrders() {
	offset := o.GetString("offset")
	limit := o.GetString("limit")
	upAgentId := o.GetString("userId")

	listRes := orderList(o.GetUserId(), offset, limit)
	response := &selectOrderListResponse{
		OrderItemMap: listRes.OrderItemMap,
		GoodsMap: listRes.GoodsMap,
	}

	selectOrderList := []*selectOrder{}
	if listRes.OrderList != nil {
		for _, order := range listRes.OrderList {
			agents := order.Agents
			sOrder := &selectOrder{
				Order: order,
			}
			if agents != nil {
				for _, agent := range agents {
					if agent.UpAgentId == upAgentId {
						sOrder.IsSelected = true
					}
				}
			}
			selectOrderList = append(selectOrderList, sOrder)
		}
	}
	response.OrderList = selectOrderList
	o.Data["json"] = response
	o.ServeJSON()
}

// @router /agent-order-list [get]
func (o *OrderController) GetAgentOrderList() {
	offset := o.GetString("offset")
	limit := o.GetString("limit")
	userId := o.GetString("userId")
	agentId := o.GetString("agentId")
	key := o.GetString("key")
	agentId = GetAgentId(agentId, key)

	o.Data["json"] = GetAgentOrders(offset, limit, userId, agentId, key)
	o.ServeJSON()
}

type confirmOrderRequest struct {
	OrderId    string     `json:"orderId" validate:"required"`
	UpdateTime *time.Time `json:"updateTime"`
}

// @router /order-confirm [post]
func (o *OrderController) UpAgentConfirmOrder() {

	var request confirmOrderRequest
	o.Bind(&request)
	userId := o.GetUserId()
	order := service.GetOrderById(request.OrderId)
	reqUpdateTime := ""
	if request.UpdateTime != nil {
		reqUpdateTime = request.UpdateTime.Format("2006-01-02 15:04:05")
	}

	if (order.UpdateTime != nil && order.UpdateTime.Format("2006-01-02 15:04:05") == reqUpdateTime) || order.UpdateTime == nil {
		if order.Agents != nil {
			hasRole := false
			for _, agent := range order.Agents {
				if agent.UpAgentId == userId {
					hasRole = true
					break
				}
			}
			if hasRole {
				order.OwnerId = userId
				order = service.UpdateOrder(userId, order)
			}
		}

	}

	o.Data["json"] = order
	o.ServeJSON()
}

func GetAgentOrders(offset string, limit string, userId string, agentId string, key string) *OrderListResponse {
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)
	agentId = GetAgentId(agentId, key)

	response := &OrderListResponse{}

	if agentId == "" {
		return response
	}

	orders := service.GetAgentOrders(agentId, userId, offsetInt, limitInt)
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

	return response
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
	utils.WriteXlsx(o.Ctx, getOrderRecords(orders, orderItemMap, goodsMap), fileName)
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
		"日期",
		"姓名",
		"电话号码",
		"地址",
		"商品",
		"总价",
		"备注",
		"快递单号",
	})
	for _, order := range orders {
		record := []string{}
		record = append(record, order.CreatedTime.Format("2006-01-02"))
		record = append(record, order.Name)
		record = append(record, order.Tel)
		record = append(record, order.Address)

		items := orderItemMap[order.Id]
		goods := ""
		for _, item := range items {
			goods += goodsMap[item.GoodsId].Name + "x" + strconv.Itoa(item.Quantity) + ", "
		}
		record = append(record, goods)
		record = append(record, strconv.Itoa(order.TotalPrice))
		record = append(record, order.Note)
		record = append(record, order.Express)

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