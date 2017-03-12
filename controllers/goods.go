package controllers

import (
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/models"
	"strconv"
)

type GoodsController struct {
	BaseController
}

type goodsRequest struct {
	Id      string     `json:"id"`
	Name    string     `json:"name" validate:"required"`
}

// @router /merge [post]
func (c *GoodsController) MergeGoods() {
	var request goodsRequest
	c.Bind(&request)

	goods := &models.Goods{
		Name: request.Name,
		UserId: c.GetUserId(),
	}
	if request.Id == "" {
		goods = service.AddGoods(goods)
	} else {
		goods.Id = request.Id
		goods = service.UpdateGoods(goods)
	}
	c.Data["json"] = goods
	c.ServeJSON()
}

type goodsRemoveRequest struct {
	Id string     `json:"id"`
}
// @router /remove [post]
func (c *GoodsController) RemoveGoods() {
	var request goodsRemoveRequest
	c.Bind(&request)

	service.RemoveGoodsById(request.Id)
	c.Data["json"] = map[string]bool{"success": true}
	c.ServeJSON()
}

// @router /list [get]
func (c *GoodsController) GetGoods() {
	offset := c.GetString("offset")
	limit := c.GetString("limit")
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)

	goods := service.GetGoods(c.GetUserId(), offsetInt, limitInt)
	c.Data["json"] = goods
	c.ServeJSON()
}