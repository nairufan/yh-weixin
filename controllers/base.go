package controllers

import (
	"github.com/astaxie/beego"
	"gopkg.in/bluesuncorp/validator.v5"
	"encoding/json"
)

const (
	userID = "userId"
)

var validate = validator.New("validate", validator.BakedInValidators)

type BaseController struct {
	beego.Controller
}

func (b *BaseController) SetUserId(id string) {
	b.SetSession(userID, id)
}

func (b *BaseController) GetUserId() string {
	uid := b.GetSession(userID)
	if uid == nil {
		return ""
	}
	return uid.(string)
}

func (b *BaseController)Bind(request interface{}) {
	if err := json.Unmarshal(b.Ctx.Input.RequestBody, request); err != nil {
		panic(err)
	}
	errs := validate.Struct(request)
	if errs != nil {
		panic(errs)
	}
}