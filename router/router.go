package router

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/labstack/echo"
	"fmt"
	"github.com/labstack/echo/middleware"
	"github.com/beewit/pay/handler"
	"github.com/beewit/pay/global"
)

func Start() {
	fmt.Printf("登陆授权系统启动")

	e := echo.New()

	e.Use(middleware.RequestID())

	e.Static("/app", "app")
	e.File("/", "app/page/index.html")

	e.POST("/member/type", handler.GetMemberTypeAndCharge)

	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))
}
