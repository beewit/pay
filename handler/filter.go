package handler

import (
	"encoding/json"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
)

func Filter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tc, _ := c.Cookie("token")
		var token string
		if tc == nil || tc.Value == "" {
			token = c.FormValue("token")
		} else {
			token = tc.Value
		}
		if token == "" {
			return utils.AuthFail(c, "登陆信息token无效，请重新登陆")
		}

		accMapStr, err := global.RD.GetString(token)
		if err != nil {
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		accMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(accMapStr), &accMap)
		if err != nil {
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆，ERR："+err.Error())
		}
		m, err := global.DB.Query("SELECT id,nickname,photo,mobile,status FROM account WHERE id=? LIMIT 1", accMap["id"])
		if err != nil {
			return utils.AuthFail(c, "获取用户信息失败")
		}
		if len(m) <= 0 {
			return utils.AuthFail(c, "用户信息不存在")
		}
		if convert.ToString(m[0]["status"]) != enum.NORMAL {
			return utils.AuthFail(c, "用户已被冻结")
		}
		c.Set("account", global.ToMapAccount(m[0]))
		return next(c)
	}
}

func GetAccount(c echo.Context) (acc *global.Account, err error) {
	itf := c.Get("account")
	if itf == nil {
		err = utils.AuthFailNull(c)
		return
	}
	acc = global.ToInterfaceAccount(itf)
	if acc == nil {
		err = utils.AuthFailNull(c)
		return
	}
	return
}
