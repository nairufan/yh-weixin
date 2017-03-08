package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"],
		beego.ControllerComments{
			Method: "MergeCustomer",
			Router: `/merge`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"],
		beego.ControllerComments{
			Method: "RemoveCustomer",
			Router: `/remove`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:CustomerController"],
		beego.ControllerComments{
			Method: "GetCustomers",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "MergeCustomer",
			Router: `/merge`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "RemoveGoods",
			Router: `/remove`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:GoodsController"],
		beego.ControllerComments{
			Method: "GetGoods",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"],
		beego.ControllerComments{
			Method: "CreateOrder",
			Router: `/create`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"],
		beego.ControllerComments{
			Method: "UpdateOrder",
			Router: `/update`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"],
		beego.ControllerComments{
			Method: "DeleteOrderItem",
			Router: `/remove-item`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"],
		beego.ControllerComments{
			Method: "UpdateOrderItem",
			Router: `/merge-item`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:OrderController"],
		beego.ControllerComments{
			Method: "GetOrders",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"],
		beego.ControllerComments{
			Method: "WxExchangeCode",
			Router: `/wx-login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"],
		beego.ControllerComments{
			Method: "MockLogin",
			Router: `/mock-login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/nairufan/yh-weixin/controllers:UserController"],
		beego.ControllerComments{
			Method: "BuyHistory",
			Router: `/buy-history`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
