package controllers

type UserController struct {
	BaseController
}

// @router /mock-login [get]
func (u *UserController) MockLogin() {
	u.SetUserId("1111")
}