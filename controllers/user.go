package controllers

import (
	"strconv"
	"time"
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/models"
	"github.com/nairufan/yh-weixin/agent"
)

type UserController struct {
	BaseController
}

type exChangeResponse struct {
	Openid     string       `json:"openid"`
	SessionKey string       `json:"session_key"`
	ErrCode    int          `json:"errcode"`
	ErrMsg     string       `json:"errmsg"`
}

type loginResponse struct {
	SessionId string     `json:"sessionId"`
}

// @router /wx-login [get]
func (u *UserController) WxExchangeCode() {
	code := u.GetString("code")

	token := agent.MustGetXCXAccessToken(code)
	user := service.GetUserByOpenId(token.OpenId)

	if user == nil {
		user = service.AddUser(&models.User{
			OpenId: token.OpenId,
			UnionId: token.UnionId,
		})
		initUserDefaultData(user.Id)
	} else if user.UnionId == "" {
		user.UnionId = token.UnionId
		user = service.UpdateUser(user)
	}

	u.SetUserId(user.Id)

	response := &loginResponse{
		SessionId: u.CruSession.SessionID(),
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @router /wx-qrc-login [get]
func (u *UserController) WxQRCLogin() {
	code := u.GetString("code")
	token := agent.MustGetQRAccessToken(code)
	user := service.GetUserByUnionId(token.UnionId)
	if user == nil {
		u.Redirect("/wx/404", 301)
		return;
	}
	u.SetUserId(user.Id)
	u.Redirect("/wx/html/download", 301)
}

// @router /mock-login [get]
func (u *UserController) MockLogin() {
	id := u.GetString("id")
	if id == "" {
		id = "111"
	}
	u.SetUserId(id)
	response := &loginResponse{
		SessionId: u.CruSession.SessionID(),
	}
	u.Data["json"] = response
	u.ServeJSON()
}

type goodsType struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Quantity int                    `json:"quantity"`
}

type buyHistory struct {
	OrderId     string                 `json:"orderId"`
	GoodsList   []*goodsType           `json:"goodsList"`
	CreatedTime *time.Time             `json:"createdTime"`
}

// @router /buy-history [get]
func (u *UserController) BuyHistory() {
	tel := u.GetString("tel")
	offset := u.GetString("offset")
	limit := u.GetString("limit")
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)
	response := []*buyHistory{}

	orders := service.GetOrderByTel(u.GetUserId(), tel, offsetInt, limitInt)
	if len(orders) > 0 {
		orderIds := []string{}
		for _, order := range orders {
			orderIds = append(orderIds, order.Id)
			response = append(response, &buyHistory{
				OrderId: order.Id,
				CreatedTime: order.CreatedTime,
			})
		}
		orderGoodsMap := map[string][]*goodsType{}
		orderItems := service.GetOrderItems(orderIds)
		goodsIds := []string{}
		for _, orderItem := range orderItems {
			goodsIds = append(goodsIds, orderItem.GoodsId)
			goodsL := orderGoodsMap[orderItem.OrderId]
			if goodsL == nil {
				goodsL = []*goodsType{}
			}
			goodsL = append(goodsL, &goodsType{
				Id: orderItem.GoodsId,
				Quantity: orderItem.Quantity,
			})
			orderGoodsMap[orderItem.OrderId] = goodsL
		}
		goodsList := service.GetGoodsByIds(goodsIds)
		goodsMap := ConvertGoodsMap(goodsList)

		for _, g := range response {
			g.GoodsList = orderGoodsMap[g.OrderId]
			for _, goods := range g.GoodsList {
				goods.Name = goodsMap[goods.Id].Name
			}
		}
	}

	u.Data["json"] = response
	u.ServeJSON()
}

// @router /statistics [get]
func (u *UserController) UserStatistics() {
	response := models.StatisticResponse{}
	now := time.Now()
	start := now.AddDate(0, 0, -10)
	statistics := service.UserStatistics(start, now)
	total := service.UserCount()
	response.Statistics = statistics
	response.Total = total
	u.Data["json"] = response
	u.ServeJSON()
}

func initUserDefaultData(userId string) {

	customerA := &models.Customer{
		Name: "范冰冰",
		Tel: "11223344",
		Address: "中国山东",
		UserId: userId,
	}
	customerB := &models.Customer{
		Name: "李易峰",
		Tel: "11223355",
		Address: "中国四川",
		UserId: userId,
	}
	customerC := &models.Customer{
		Name: "吴亦凡",
		Tel: "11223366",
		Address: "中国北京",
		UserId: userId,
	}

	customerA = service.AddCustomer(customerA)
	customerB = service.AddCustomer(customerB)
	customerC = service.AddCustomer(customerC)

	goodsA := &models.Goods{
		Name: "小番茄",
		UserId: userId,
	}

	goodsB := &models.Goods{
		Name: "西瓜",
		UserId: userId,
	}

	goodsA = service.AddGoods(goodsA)
	goodsB = service.AddGoods(goodsB)

	orderA := initAddOrder(userId, customerA, []*models.Goods{
		goodsA, goodsB,
	}, 3)
	initAddOrder(userId, customerB, []*models.Goods{
		goodsB,
	}, 1)
	orderC := initAddOrder(userId, customerC, []*models.Goods{
		goodsA,
	}, 6)

	orderA.Status = models.OrderStatusDone
	orderC.Status = models.OrderStatusDone
	service.UpdateOrder(orderA)
	service.UpdateOrder(orderC)
}

func initAddOrder(userId string, customer *models.Customer, goods []*models.Goods, quantity int) *models.Order {
	order := &models.Order{
		UserId: userId,
		Name: customer.Name,
		Tel: customer.Tel,
		Address: customer.Address,
		TotalPrice: 888,
		Note: "尽快发货",
	}
	order = service.AddOrder(order)
	for _, goods := range goods {
		orderItem := &models.OrderItem{
			OrderId: order.Id,
			GoodsId: goods.Id,
			Quantity: quantity,
		}
		orderItem = service.AddOrderItem(orderItem)
	}

	return order
}