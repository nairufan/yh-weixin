package main

import (
	_ "github.com/nairufan/yh-weixin/docs"
	_ "github.com/nairufan/yh-weixin/routers"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego"
	"github.com/nairufan/yh-weixin/filters"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.BConfig.WebConfig.StaticDir["/html/js"] = "static/js"

	// api filter
	beego.InsertFilter("/api/customer/merge", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/customer/remove", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/customer/list", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/goods/merge", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/goods/remove", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/goods/list", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/create", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/update", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/update-items", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/remove-item", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/merge-item", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/list", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/user/buy-history", beego.BeforeRouter, filters.LoginCheck)

	//html page
	beego.InsertFilter("/html/download", beego.BeforeRouter, filters.LoginCheck)

	beego.Run()
}
