package filters

import (
	"github.com/astaxie/beego/context"
)

const (
	userID = "userId"
	Role = "role"
)

func LoginCheck(ctx *context.Context) {
	_, ok := ctx.Input.Session(userID).(string)
	if !ok {
		panic("403")
	}
}
