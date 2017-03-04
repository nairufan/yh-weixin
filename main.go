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
	// api filter
	beego.InsertFilter("/api/customer/*", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/goods/*", beego.BeforeRouter, filters.LoginCheck)
	beego.InsertFilter("/api/order/*", beego.BeforeRouter, filters.LoginCheck)
	beego.Run()
}
