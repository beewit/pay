package handler

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"fmt"

	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/pay/alipay"
	"github.com/beewit/pay/global"
	"github.com/beewit/pay/wxpay"
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

//多个功能开通
func CreateBatchFuncOrder(c echo.Context) error {
	funcIdStr := c.FormValue("funcId")
	fcIdStr := c.FormValue("fcId")
	pt := c.FormValue("pt")
	if funcIdStr == "" {
		return utils.Error(c, "请正确选择开通功能", nil)
	}
	if fcIdStr == "" || !utils.IsValidNumber(fcIdStr) {
		return utils.Error(c, "请正确选择功能开通", nil)
	}
	fcId, _ := strconv.ParseInt(fcIdStr, 10, 64)
	mt := getFuncList(funcIdStr)
	fc := getFuncCharge(fcId)
	if mt == nil {
		return utils.Error(c, "选择的开通功能不存在", nil)
	}
	if fc == nil {
		return utils.Error(c, "选择的功能开通方式不存在", nil)
	}
	acc := global.ToInterfaceAccount(c.Get("account"))
	accID := acc.ID
	accName := acc.Nickname
	flog, tradeNo, codeUrl, getUrl := getOrderCode(mt, fc, accID, accName, pt, c.RealIP())
	if flog {
		return utils.Success(c, "创建订单成功", map[string]interface{}{"codeUrl": codeUrl, "getUrl": getUrl, "tradeNo": tradeNo})
	}
	return utils.Error(c, "创建订单失败", nil)
}

//单个功能开通
func CreateFuncOrder(c echo.Context) error {
	funcIdStr := c.FormValue("funcId")
	fcIdStr := c.FormValue("fcId")
	pt := c.FormValue("pt")
	if funcIdStr == "" || !utils.IsValidNumber(funcIdStr) {
		return utils.Error(c, "请正确选择开通功能", nil)
	}
	if fcIdStr == "" || !utils.IsValidNumber(fcIdStr) {
		return utils.Error(c, "请正确选择功能开通", nil)
	}
	funcId, _ := strconv.ParseInt(funcIdStr, 10, 64)
	fcId, _ := strconv.ParseInt(fcIdStr, 10, 64)
	mt := getFunc(funcId)
	fc := getFuncCharge(fcId)
	if mt == nil {
		return utils.Error(c, "选择的开通功能不存在", nil)
	}
	if fc == nil {
		return utils.Error(c, "选择的功能开通不存在", nil)
	}
	acc := global.ToInterfaceAccount(c.Get("account"))
	accID := acc.ID
	accName := acc.Nickname
	flog, tradeNo, codeUrl, getUrl := getOrderCode(mt, fc, accID, accName, pt, c.RealIP())
	if flog {
		return utils.Success(c, "创建订单成功", map[string]interface{}{"codeUrl": codeUrl, "getUrl": getUrl, "tradeNo": tradeNo})
	}
	return utils.Error(c, "创建订单失败", nil)
}

func GetFuncAndCharge(c echo.Context) error {
	fid := c.FormValue("fid")
	if fid == "" {
		return utils.ErrorNull(c, "请选择开通功能")
	}
	m := make(map[string]interface{})
	m["account"] = global.ToInterfaceAccount(c.Get("account"))
	sql := fmt.Sprintf("SELECT * FROM func WHERE status=? AND id in(%s) ORDER BY `order` DESC,ct_time DESC", fid)
	rows, err := global.DB.Query(sql, enum.NORMAL)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "获取开通功能异常", nil)
	}
	if len(rows) < 1 {
		return utils.NullData(c)
	}
	mt := rows
	m["func"] = mt
	sql = "SELECT * FROM func_charge WHERE status=? ORDER BY `order` DESC,ct_time DESC"
	rows2, err2 := global.DB.Query(sql, enum.NORMAL)
	if err2 != nil {
		global.Log.Error(err2.Error())
		return utils.Error(c, "获取功能开通异常", nil)
	}
	if len(rows2) > 0 {
		fc := rows2
		m["funcCharge"] = fc
	}
	return utils.Success(c, "", m)
}

func AlipayNotify(c echo.Context) error {
	args, _ := c.FormParams()
	query, _ := url.QueryUnescape(args.Encode())
	global.Log.Info("支付成功的参数：" + query)
	err := alipay.NewTrade().Verify(args, global.AlipayAliPublicPath)
	if err != nil {
		global.Log.Error(err.Error())
		return c.HTML(http.StatusOK, "error")
	}
	UpdateOrderFuncStatus(convert.MustInt64(c.FormValue("out_trade_no")), convert.MustFloat64(c.FormValue("total_amount")))
	return c.HTML(http.StatusOK, "success")
}

func WechatNotify(c echo.Context) error {
	body, bErr := ioutil.ReadAll(c.Request().Body)
	if bErr != nil {
		global.Log.Error("读取http body失败，原因：", bErr.Error())
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  "读取http body失败",
		}
		return c.XML(http.StatusOK, wr)
	}
	defer c.Request().Body.Close()

	var args wxpay.Notice
	bErr = xml.Unmarshal(body, &args)
	if bErr != nil {
		global.Log.Error("读取http body失败，原因：", bErr.Error())
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  "解析HTTP Body格式到xml失败",
		}
		return c.XML(http.StatusOK, wr)
	}

	//args, _ := c.FormParams()
	j, _ := json.Marshal(args)
	global.Log.Info("支付成功的参数：" + string(j[:]))
	err := wxpay.NewTrade().Verify(args, global.WechatApiKey)
	if err != nil {
		global.Log.Error(err.Error())
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  "验签失败",
		}
		return c.XML(http.StatusOK, wr)
	}
	UpdateOrderFuncStatus(convert.MustInt64(args.OutTradeNo), convert.MustFloat64(args.TotalFee)/100)
	wr := &wxpay.Response{
		ReturnCode: "SUCCESS",
	}
	return c.XML(http.StatusOK, wr)
}

func GetOrderById(c echo.Context) error {
	tradeNo := c.FormValue("tradeNo")
	if tradeNo == "" || !utils.IsValidNumber(tradeNo) {
		return utils.Error(c, "无有效订单", nil)
	}
	order := getOrderFuncById(convert.MustInt64(tradeNo))
	if order == nil {
		return utils.NullData(c)
	}
	if order["pay_status"] == enum.PAY_STATUS_END {
		return utils.Success(c, "已支付", order)
	}
	return utils.NullData(c)
}

func getFunc(id int64) []map[string]interface{} {
	sql := "SELECT * FROM func WHERE status=? AND id=? LIMIT 1"
	rows, err := global.DB.Query(sql, enum.NORMAL, id)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows
}

func getFuncList(ids string) []map[string]interface{} {
	sql := fmt.Sprintf("SELECT * FROM func WHERE status=? AND id in(%s)", ids)
	rows, err := global.DB.Query(sql, enum.NORMAL)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows
}

func getFuncCharge(id int64) map[string]interface{} {
	sql := "SELECT * FROM func_charge WHERE status=? AND id=? LIMIT 1"
	rows, err := global.DB.Query(sql, enum.NORMAL, id)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func updateOrderUrl(id int64, codeUrl, getUrl string) int64 {
	sql := "UPDATE order_payment SET code_url=?,get_url=? WHERE id=?"
	x, err := global.DB.Update(sql, codeUrl, getUrl, id)
	if err != nil {
		global.Log.Error(err.Error())
		return 0
	}
	if x < 1 {
		return 0
	}
	return x
}

func UpdateOrderFuncStatus(id int64, price float64) bool {
	flog := false
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		errMsg := ""
		ids := convert.ToString(id)
		//查询订单
		order := getOrderFuncById(id)
		if order == nil {
			errMsg := "订单不存在：" + ids
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}
		if order["pay_status"] == enum.PAY_STATUS_END {
			errMsg := "订单已支付过了：" + ids
			global.Log.Error(errMsg)
			panic(errMsg)
		}
		if convert.MustFloat64(order["price"]) != price {
			errMsg := "订单支付金额不一致"
			global.Log.Warning(errMsg)
			panic(errMsg)
			return
		}
		var x int64
		var err error
		funcId, errFuncId := convert.ToInt64(order["func_id"])
		if errFuncId != nil {
			global.Log.Error(ids + "订单付款的开通功能异常：" + errFuncId.Error())
			panic(err)
		}

		accId, errAccId := convert.ToInt64(order["account_id"])
		if errAccId != nil {
			global.Log.Error(ids + "订单付款的帐号异常：" + errAccId.Error())
			panic(err)
		}

		days, errDays := convert.ToInt(order["days"])
		if errDays != nil {
			global.Log.Error(ids + "订单付款-功能开通：" + errDays.Error())
			panic(err)
		}
		sql := "UPDATE order_payment SET pay_status=?,pay_time=? WHERE id=?"
		x, err = tx.Update(sql, enum.PAY_STATUS_END, utils.CurrentTime(), id)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			errMsg := "修改订单状态失败：" + ids
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}

		var daysTime time.Time
		accFunc := getAccountFuncId(accId, funcId)
		if accFunc != nil {
			//新增用户开通功能清单数据
			expirTimeStr, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(accFunc["expiration_time"]))
			if errExpirTime != nil {
				global.Log.Error(convert.ToString(accId) + "会员的过期时间错误：" + errExpirTime.Error())
				panic(err)
			}
			if expirTimeStr.After(time.Now()) {
				//未到期的续费
				daysTime = expirTimeStr.AddDate(0, 0, days)
			} else {
				//已到期的续费
				daysTime = time.Now().AddDate(0, 0, days)
			}
			sql = "UPDATE account_func SET expiration_time=?,ut_time=? WHERE account_id=? AND func_id=?"
			x, err = tx.Update(sql, utils.FormatTime(daysTime), utils.CurrentTime(), accId, funcId)
		} else {
			//修改已经开通过的清单数据
			//个人升级为企业，或第一次开通
			daysTime = time.Now().AddDate(0, 0, days)
			iw, _ := utils.NewIdWorker(1)
			id, _ := iw.NextId()
			m := make(map[string]interface{})
			m["id"] = id
			m["account_id"] = accId
			m["func_id"] = funcId
			m["expiration_time"] = utils.FormatTime(daysTime)
			m["ct_time"] = utils.CurrentTime()
			m["ut_time"] = utils.CurrentTime()
			_, err = tx.InsertMap("account_func", m)
		}
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		flog = true
		errMsg = "修改订单成功：" + ids
		global.Log.Info(errMsg)
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
			flog = false
		}
	})
	return flog
}

func getAccountId(id int64) map[string]interface{} {
	sql := "SELECT * FROM account WHERE status=? AND id=? LIMIT 1"
	rows, err := global.DB.Query(sql, enum.NORMAL, id)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getAccountFuncId(accId, funcId int64) map[string]interface{} {
	sql := "SELECT * FROM account_func WHERE account_id=? AND func_id=? LIMIT 1"
	rows, err := global.DB.Query(sql, accId, funcId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getOrderFuncById(id int64) map[string]interface{} {
	sql := "SELECT o.*,om.func_id,om.func_name,om.func_charge_id,om.days FROM order_payment o LEFT JOIN " +
		"order_payment_record_func om ON o.id=om.order_payment_id WHERE o.id=? LIMIT 1"
	rows, err := global.DB.Query(sql, id)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getOrderCode(mt []map[string]interface{}, fc map[string]interface{}, accId int64, accName, payType, ip string) (bool, int64, string,
	string) {
	m := make(map[string]interface{})
	iw, _ := utils.NewIdWorker(1)
	tradeNo, _ := iw.NextId()

	var sumPrice float64 = 0

	funcList := make([]map[string]interface{}, len(mt))
	for i := 0; i < len(mt); i++ {
		mr := make(map[string]interface{})
		mrId, _ := iw.NextId()
		mr["id"] = mrId
		mr["order_payment_id"] = tradeNo
		mr["func_id"] = mt[i]["id"]
		mr["func_name"] = mt[i]["name"]
		mr["price"] = mt[i]["price"]
		mr["func_charge_id"] = fc["id"]
		mr["days"] = fc["days"]
		funcList[i] = mr

		sumPrice += convert.MustFloat64(mr["price"])
	}

	totalPrice := convert.MustFloat64(fc["days"]) * sumPrice
	//测试使用
	totalPrice = 0.01

	m["id"] = tradeNo
	m["account_id"] = accId
	m["account_name"] = accName
	m["type"] = enum.PAY_TYPE_FUNC
	m["pay_type"] = payType
	m["pay_price"] = totalPrice
	m["pay_status"] = enum.PAY_STATUS_NOT
	m["status"] = enum.NORMAL
	m["ct_time"] = utils.CurrentTime()
	m["ct_ip"] = ip

	flog := true
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		_, err := tx.InsertMap("order_payment", m)
		if err != nil {
			panic(err)
		}
		_, err = tx.InsertMapList("order_payment_record_func", funcList)
		if err != nil {
			panic(err)
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
			flog = false
		}
	})

	if flog {
		var codeUrl, getUrl string
		var payErr error
		if payType == enum.PAY_TYPE_ALIPAY {
			codeUrl, getUrl, payErr = alipay.GetPayUrl(
				"工蜂小智 - 功能开通",
				"工蜂小智 - 功能开通",
				convert.ToString(tradeNo),
				totalPrice)
		} else if payType == enum.PAY_TYPE_WECHAT {
			codeUrl, payErr = wxpay.GetPayUrl(
				"工蜂小智 - 功能开通",
				"工蜂小智 - 功能开通",
				convert.ToString(tradeNo),
				totalPrice)
		}
		if payErr != nil {
			return false, 0, "", ""
		}
		updateOrderUrl(tradeNo, codeUrl, getUrl)
		return flog, tradeNo, codeUrl, getUrl
	}
	return false, 0, "", ""
}
