package controllers

import (
	"github.com/nairufan/yh-weixin/service"
	"github.com/nairufan/yh-weixin/models"
	"strconv"
	"time"
)

type CustomerController struct {
	BaseController
}

type customerRequest struct {
	Id      string     `json:"id"`
	Name    string     `json:"name" validate:"required"`
	Tel     string     `json:"tel" validate:"required"`
	Address string     `json:"address"`
	Note    string     `json:"note"`
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

	customer.Name = request.Name
	customer.Tel = request.Tel
	customer.Address = request.Address
	customer.Note = request.Note

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

// @router /statistics [get]
func (c *CustomerController) CustomerStatistics() {
	response := models.StatisticResponse{}
	now := time.Now()
	start := now.AddDate(0, 0, -10)
	statistics := service.CustomerStatistics(start, now)
	total := service.CustomerCount()
	response.Statistics = statistics
	response.Total = total
	c.Data["json"] = response
	c.ServeJSON()
}

// @router /init_py [get]
func (c *CustomerController) CustomerPY() {
	customers := service.GetAllCustomers()
	for _, c := range customers {
		service.UpdateCustomer(c)
	}
	c.Data["json"] = customers
	c.ServeJSON()
}