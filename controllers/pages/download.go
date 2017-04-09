package pages

import (
	"io/ioutil"
	"github.com/astaxie/beego/context"
)

func Download(ctx *context.Context) {
	content, err := ioutil.ReadFile("static/uhexceldownload.html")
	if err != nil {
		panic(err)
	}
	ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Output.Body(content)
}

func Qrcode(ctx *context.Context) {
	content, err := ioutil.ReadFile("static/uhqrcode.html")
	if err != nil {
		panic(err)
	}
	ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Output.Body(content)
}