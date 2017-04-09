package filters

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego"
)

const (
	userID = "userId"
	Role = "role"
)

func LoginCheck(ctx *context.Context) {
	_, ok := ctx.Input.Session(userID).(string)
	beego.Info(ctx.Input.Session(userID))
	if !ok {
		panic("403")
	}
}
