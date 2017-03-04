package controllers

type UserController struct {
	BaseController
}

type loginResponse struct {
	SessionId string     `json:"sessionId"`
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