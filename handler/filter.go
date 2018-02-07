package handler

import (
	"encoding/json"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"strings"
	//"errors"
	//"github.com/beewit/wechat/mp"
	"github.com/beewit/wechat/mp/user/oauth2"
	"github.com/beewit/wechat/mp"
	"errors"
)

var MPSessionId string = "mpSessionId"

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

func GetMiniAppSession(c echo.Context) (*mp.WxSesstion, error) {
	miniAppSessionId := strings.TrimSpace(c.FormValue("miniAppSessionId"))
	mpSessionId := strings.TrimSpace(c.FormValue("mpSessionId"))
	if miniAppSessionId == "" && mpSessionId == "" {
		return nil, errors.New("未识别到用户标识")
	}
	var wsStr string
	var err error
	if miniAppSessionId != "" {
		wsStr, err = global.RD.GetString(miniAppSessionId)
	} else {
		wsStr, err = global.RD.GetString(mpSessionId)
	}
	if err != nil {
		return nil, errors.New("未识别到用户标识")
	}
	var ws *mp.WxSesstion
	err = json.Unmarshal([]byte(wsStr), &ws)
	if err != nil {
		return nil, errors.New("获取用户登录标识失败")
	}
	return ws, nil
}


func GetOauthUser(c echo.Context) *oauth2.UserInfo {
	mpSessionId := strings.TrimSpace(c.FormValue(MPSessionId))
	if mpSessionId == "" {
		return nil
	}
	us, err := global.RD.GetString(mpSessionId)
	if err != nil {
		return nil
	}
	var u *oauth2.UserInfo
	err = json.Unmarshal([]byte(us), &u)
	if err != nil {
		return nil
	}
	return u
}