package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
)

type UserController struct {
	BaseController
}

type exChangeResponse struct {
	Openid     string       `json:"openid"`
	SessionKey string       `json:"session_key"`
	ErrCode    int          `json:"errcode"`
	ErrMsg     string       `json:"errmsg"`
}

type loginResponse struct {
	SessionId string     `json:"sessionId"`
}

// @router /wx-login [get]
func (u *UserController) WxExchangeCode() {
	code := u.GetString("code")
	url := beego.AppConfig.String("wx.exChangeCodeUrl")
	appId := beego.AppConfig.String("wx.appid")
	secret := beego.AppConfig.String("wx.secret")
	url = fmt.Sprintf(url, appId, secret, code)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	beego.Info(string(body))
	exResponse := &exChangeResponse{}
	json.Unmarshal(body, exResponse)
	if exResponse.ErrMsg != "" {
		beego.Error(exResponse.ErrCode)
		beego.Error(exResponse.ErrMsg)
		panic(exResponse.ErrMsg)
	}
	u.SetUserId(exResponse.Openid)

	response := &loginResponse{
		SessionId: u.CruSession.SessionID(),
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @router /mock-login [get]
func (u *UserController) MockLogin() {
	u.SetUserId("1111")
	response := &loginResponse{
		SessionId: u.CruSession.SessionID(),
	}
	u.Data["json"] = response
	u.ServeJSON()
}