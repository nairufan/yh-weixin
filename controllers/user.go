package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/nairufan/yh-weixin/service"
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
	url := beego.AppConfig.String("wx.exChangeCodeUrl")
	appId := beego.AppConfig.String("wx.appid")
	secret := beego.AppConfig.String("wx.secret")
	url = fmt.Sprintf(url, appId, secret, code)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	beego.Info(string(body))
	exResponse := &exChangeResponse{}
	json.Unmarshal(body, exResponse)
	if exResponse.ErrMsg != "" {
		beego.Error(exResponse.ErrCode)
		beego.Error(exResponse.ErrMsg)
		panic(exResponse.ErrMsg)
	}
	u.SetUserId(exResponse.Openid)

	response := &loginResponse{
		SessionId: u.CruSession.SessionID(),
	}
	u.Data["json"] = response
	u.ServeJSON()
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