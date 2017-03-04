// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/nairufan/yh-weixin/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/customer",
			beego.NSInclude(
				&controllers.CustomerController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/goods",
			beego.NSInclude(
				&controllers.GoodsController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&controllers.OrderController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
