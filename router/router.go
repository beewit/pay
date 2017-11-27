package router

import (
	"fmt"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/pay/global"
	"github.com/beewit/pay/handler"
	"github.com/labstack/echo"
)

func Start() {
	fmt.Printf("登陆授权系统启动")

	e := echo.New()
	e.Static("/app", "app")
	e.File("/", "app/page/index.html")
	e.File("favicon.ico", "app/static/img/favicon.ico")
	e.POST("/order/query", handler.GetOrderById, handler.Filter)
	e.POST("/order/create", handler.CreateFuncOrder, handler.Filter)
	e.POST("/order/create/list", handler.CreateBatchFuncOrder, handler.Filter)
	e.POST("/order/app/create", handler.CreateAppOrder, handler.Filter)
	e.POST("/member/type", handler.GetFuncAndCharge, handler.Filter)
	e.POST("/alipay/notify", handler.AlipayNotify)
	e.POST("/wechat/notify", handler.WechatNotify)

	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))

}
