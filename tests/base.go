package test

import (
	"net/http"
	"net/http/httptest"
	"github.com/astaxie/beego"
	"encoding/json"
	"time"
	"bytes"
	"math/rand"
)

type loginResponse struct {
	SessionId string     `json:"sessionId"`
}

func Login() string {
	r, _ := http.NewRequest("GET", "/api/user/mock-login", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	loginResponse := &loginResponse{}
	json.Unmarshal(w.Body.Bytes(), loginResponse)

	return loginResponse.SessionId
}

func GenRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func DoRequest(input interface{}, output interface{}, sessionId string) {
	b, _ := json.Marshal(input)
	r, _ := http.NewRequest("POST", "/api/goods/merge", bytes.NewReader(b))
	r.AddCookie(&http.Cookie{
		Name: "sessionId",
		Value: sessionId,
		Path: "/",
	})
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	json.Unmarshal(w.Body.Bytes(), output)
}