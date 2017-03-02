package controllers

import (
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/models"
	"strconv"
)

type CustomerController struct {
	BaseController
}

type customerRequest struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Tel     string     `json:"tel"`
	Address string     `json:"address"`
}

// @router /merge [post]
func (c *CustomerController) MergeCustomer() {
	var request customerRequest
	c.Bind(&request)

	customer := &models.Customer{
		Name: request.Name,
		Tel: request.Tel,
		Address: request.Address,
		UserId: c.GetUserId(),
	}

	if request.Name != "" {
		customer.Name = request.Name
	}
	if request.Tel != "" {
		customer.Tel = request.Tel
	}
	if request.Address != "" {
		customer.Address = request.Address
	}

	if request.Id == "" {
		customer = service.AddCustomer(customer)
	} else {
		customer.Id = request.Id
		customer = service.UpdateCustomer(customer)
	}
	c.Data["json"] = customer
	c.ServeJSON()
}

type customerRemoveRequest struct {
	Id string     `json:"id"`
}
// @router /remove [post]
func (c *CustomerController) RemoveCustomer() {
	var request customerRemoveRequest
	c.Bind(&request)

	service.RemoveCustomerById(request.Id)
	c.Data["json"] = map[string]bool{"success": true}
	c.ServeJSON()
}

// @router /list [get]
func (c *CustomerController) GetCustomers() {
	offset := c.GetString("offset")
	limit := c.GetString("limit")
	offsetInt, _ := strconv.Atoi(offset)
	limitInt, _ := strconv.Atoi(limit)

	goods := service.GetCustomers(c.GetUserId(), offsetInt, limitInt)
	c.Data["json"] = goods
	c.ServeJSON()
}