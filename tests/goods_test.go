package test

import (
	"testing"
	"runtime"
	"path/filepath"
	_ "github.com/nairufan/yh-weixin/routers"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/nairufan/yh-weixin/models"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGoods(t *testing.T) {

	sessionId := Login()

	Convey("Test Create Goods\n", t, func() {
		goodsRequest := &models.Goods{
			Name: "TestUser_" + GenRandomString(5),
		}
		goods := &models.Goods{}

		DoRequest(goodsRequest, goods, sessionId)

		So(goods.Id, ShouldNotBeNil)
		So(goods.Name, ShouldEqual, goodsRequest.Name)
	})
}