package agent

import (
	"github.com/astaxie/beego"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type exChangeResponse struct {
	Openid     string       `json:"openid"`
	SessionKey string       `json:"session_key"`
	ErrCode    int          `json:"errcode"`
	ErrMsg     string       `json:"errmsg"`
}

func MustGetOpenId(code string) string {
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
	response := &exChangeResponse{}

	json.Unmarshal(body, response)

	if response.ErrMsg != "" {
		beego.Error(response.ErrCode)
		beego.Error(response.ErrMsg)
		panic(response.ErrMsg)
	}

	return response.Openid
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	OpenId       string       `json:"openid"`
	UnionId      string       `json:"unionid"`
	Scope        string       `json:"scope"`
	ErrCode      int          `json:"errcode"`
	ErrMsg       string       `json:"errmsg"`
}

func MustGetAccessToken(code string) *TokenResponse {
	at_url := beego.AppConfig.String("wx.access_token")
	appId := beego.AppConfig.String("wx.appid")
	secret := beego.AppConfig.String("wx.secret")
	at_url = fmt.Sprintf(at_url, appId, secret, code)
	resp, err := http.Get(at_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	beego.Info(string(body))
	response := &TokenResponse{}

	json.Unmarshal(body, response)

	if response.ErrMsg != "" {
		beego.Error(response.ErrCode)
		beego.Error(response.ErrMsg)
		panic(response.ErrMsg)
	}
	return response
}

type UserResponse struct {
	OpenId     string       `json:"openid"`
	Nickname   string       `json:"nickname"`
	Sex        int       `json:"sex"`
	Province   string       `json:"province"`
	City       string       `json:"city"`
	Country    string       `json:"country"`
	HeadImgUrl string       `json:"headimgurl"`
	Privilege  string       `json:"privilege"`
	UnionId    string       `json:"unionid"`
	ErrCode    int          `json:"errcode"`
	ErrMsg     string       `json:"errmsg"`
}

func MustGetUserInfo(token string, openId string) *UserResponse {
	user_url := beego.AppConfig.String("wx.user_info")
	user_url = fmt.Sprintf(user_url, token, openId)
	beego.Info(user_url)
	resp, err := http.Get(user_url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	beego.Info(string(body))
	response := &UserResponse{}

	json.Unmarshal(body, response)

	if response.ErrMsg != "" {
		beego.Error(response.ErrCode)
		beego.Error(response.ErrMsg)
		panic(response.ErrMsg)
	}

	return response
}
