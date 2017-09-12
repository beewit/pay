package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/pay/global"
	"github.com/labstack/echo"
	"strconv"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/pay/alipay"
	"net/http"
	"net/url"
	"errors"
	"time"
	"github.com/beewit/pay/wxpay"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
)

func CreateMemberTypeOrder(c echo.Context) error {
	mtIdStr := c.FormValue("mtId")
	mtcIdStr := c.FormValue("mtcId")
	pt := c.FormValue("pt")
	if mtIdStr == "" || !utils.IsValidNumber(mtIdStr) {
		return utils.Error(c, "请正确选择会员类型", nil)
	}
	if mtcIdStr == "" || !utils.IsValidNumber(mtcIdStr) {
		return utils.Error(c, "请正确选择会员套餐", nil)
	}
	mtId, _ := strconv.ParseInt(mtIdStr, 10, 64)
	mtcId, _ := strconv.ParseInt(mtcIdStr, 10, 64)
	mt := getMemberType(mtId)
	mtc := getMemberTypeCharge(mtcId)
	if mt == nil {
		return utils.Error(c, "选择的会员类型不存在", nil)
	}
	if mtc == nil {
		return utils.Error(c, "选择的会员套餐不存在", nil)
	}
	var accId int64
	accId = 3
	accName := "曾老师"
	mtc["price"] = 0.01
	flog, tradeNo, codeUrl, getUrl := getOrderCode(mt, mtc, accId, accName, enum.PAY_TYPE_MEMBER_SET_MEAL, pt, c.RealIP())
	if flog {
		return utils.Success(c, "创建订单成功", map[string]interface{}{"codeUrl": codeUrl, "getUrl": getUrl, "tradeNo": tradeNo})
	} else {
		return utils.Error(c, "创建订单失败", nil)
	}
}

func GetMemberTypeAndCharge(c echo.Context) error {
	m := make(map[string]interface{})
	sql := "SELECT * FROM member_type WHERE status=? ORDER BY `order` DESC,ct_time DESC"
	rows, err := global.DB.Query(sql, enum.NORMAL)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.Error(c, "获取会员类型异常", nil)
	}
	if len(rows) < 1 {
		return utils.NullData(c)
	}
	mt := rows
	m["memberType"] = mt
	sql = "SELECT * FROM member_type_charge WHERE status=? ORDER BY `order` DESC,ct_time DESC"
	rows2, err2 := global.DB.Query(sql, enum.NORMAL)
	if err2 != nil {
		global.Log.Error(err2.Error())
		return utils.Error(c, "获取会员套餐异常", nil)
	}
	if len(rows2) > 0 {
		mtc := rows2
		m["memberTypeCharge"] = mtc
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
	UpdateOrderMTCStatus(convert.MustInt64(c.FormValue("out_trade_no")), convert.MustFloat64(c.FormValue("total_amount")))
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
	UpdateOrderMTCStatus(convert.MustInt64(args.OutTradeNo), convert.MustFloat64(args.TotalFee)/100)
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
	order := getOrderMTCById(convert.MustInt64(tradeNo))
	if order == nil {
		return utils.NullData(c)
	} else {
		if order["pay_status"] == enum.PAY_STATUS_END {
			return utils.Success(c, "已支付", order)
		} else {
			return utils.NullData(c)
		}
	}
}

func getMemberType(id int64) map[string]interface{} {
	sql := "SELECT * FROM member_type WHERE status=? AND id=? LIMIT 1"
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

func getMemberTypeCharge(id int64) map[string]interface{} {
	sql := "SELECT * FROM member_type_charge WHERE status=? AND id=? LIMIT 1"
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

func UpdateOrderMTCStatus(id int64, price float64) bool {
	flog := false
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		errMsg := ""
		ids := convert.ToString(id)
		//查询订单
		order := getOrderMTCById(id)
		if order == nil {
			errMsg := "订单不存在：" + ids
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}
		if order["pay_status"] == enum.PAY_STATUS_END {
			errMsg := "订单已支付过了：" + ids
			global.Log.Error(errMsg)
			return
		}
		//if convert.MustFloat64(order["price"]) != price {
		//	errMsg := "订单支付金额不一致"
		//	global.Log.Warning(errMsg)
		//	return
		//}
		sql := "UPDATE order_payment SET pay_status=?,pay_time=? WHERE id=?"
		x, err := global.DB.Update(sql, enum.PAY_STATUS_END, utils.CurrentTime(), id)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			errMsg := "修改订单状态失败：" + ids
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}
		accId, errAccId := convert.ToInt64(order["account_id"])
		if errAccId != nil {
			global.Log.Error(ids + "订单付款帐号异常：" + errAccId.Error())
		}

		mtId, errMtId := convert.ToInt64(order["member_type_id"])
		if errMtId != nil {
			global.Log.Error(ids + "订单付款会员套餐：" + errMtId.Error())
		}

		days, errDays := convert.ToInt(order["days"])
		if errDays != nil {
			global.Log.Error(ids + "订单付款会员套餐：" + errDays.Error())
		}
		acc := getAccountId(accId)
		if acc == nil {
			global.Log.Error(convert.ToString(accId) + "会员不存在")
			return
		}
		accMtId, errAccMtId := convert.ToInt64(order["member_type_id"])
		if errAccMtId != nil {
			accMtId = 0
		}

		var daysTime time.Time
		if mtId == accMtId && acc["member_expir_time"] != nil {
			//续费
			expirTimeStr, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(acc["member_expir_time"]))
			if errExpirTime != nil {
				global.Log.Error(convert.ToString(accId) + "会员的过期时间错误：" + errExpirTime.Error())
			}
			if expirTimeStr.After(time.Now()) {
				//未到期的续费
				daysTime = expirTimeStr.AddDate(0, 0, days)
			} else {
				//已到期的续费
				daysTime = time.Now().AddDate(0, 0, days)
			}
		} else {
			//个人升级为企业，或第一次开通
			daysTime = time.Now().AddDate(0, 0, days)
		}
		sql = "UPDATE account SET member_type_id=?,member_type_name=?,member_expir_time=? WHERE id=?"
		x, err = global.DB.Update(sql, order["member_type_id"], order["member_type_name"], utils.FormatTime(daysTime), accId)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			errMsg := "修改帐号会员套餐失败：" + ids
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
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

func getOrderMTCById(id int64) map[string]interface{} {
	sql := "SELECT o.*,om.member_type_id,om.member_type_name,om.member_type_charge_id,om.days FROM order_payment o LEFT JOIN " +
		"order_payment_record_mtc om ON o.id=om.order_payment_id WHERE o.id=? LIMIT 1"
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

func getEffectiveOrderMTC(accId, mtId, mtcId, days int64, t, payType string, payPrice float64) map[string]interface{} {
	sql := "SELECT o.*,om.member_type_id,om.member_type_charge_id FROM order_payment o LEFT JOIN order_payment_record_mtc om ON o.id=om.order_payment_id" +
		" WHERE o.account_id=? AND o.type=? AND o.pay_type=? AND o.pay_price=? AND o.pay_status=? AND o.status=? AND om.member_type_id=?" +
		" AND om.member_type_charge_id=? AND om.days=? AND to_days(o.ct_time) = to_days(now()) LIMIT 1"
	rows, err := global.DB.Query(sql, accId, t, payType, payPrice, enum.PAY_STATUS_NOT, enum.NORMAL, mtId, mtcId, days)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getOrderCode(mt map[string]interface{}, mtc map[string]interface{}, accId int64, accName, t, payType, ip string) (bool, int64, string,
	string) {
	order := getEffectiveOrderMTC(accId, convert.MustInt64(mt["id"]), convert.MustInt64(mtc["id"]), convert.MustInt64(mtc["days"]), t, payType,
		convert.MustFloat64(mtc["price"]))
	if order != nil && len(order) > 0 {
		return true, convert.MustInt64(order["id"]), convert.ToString(order["code_url"]), convert.ToString(order["get_url"])
	} else {
		m := make(map[string]interface{})
		iw, _ := utils.NewIdWorker(1)
		tradeNo, _ := iw.NextId()

		m["id"] = tradeNo
		m["account_id"] = accId
		m["account_name"] = accName
		m["type"] = enum.PAY_TYPE_MEMBER_SET_MEAL
		m["pay_type"] = payType
		m["pay_price"] = mtc["price"]
		m["pay_status"] = enum.PAY_STATUS_NOT
		m["status"] = enum.NORMAL
		m["ct_time"] = utils.CurrentTime()
		m["ct_ip"] = ip

		mr := make(map[string]interface{})
		mrId, _ := iw.NextId()
		mr["id"] = mrId
		mr["order_payment_id"] = tradeNo
		mr["member_type_id"] = mt["id"]
		mr["member_type_name"] = mt["name"]
		mr["member_type_charge_id"] = mtc["id"]
		mr["price"] = mtc["price"]
		mr["days"] = mtc["days"]

		flog := true
		global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
			_, err := global.DB.InsertMap("order_payment", m)
			if err != nil {
				panic(err)
			}
			_, err = global.DB.InsertMap("order_payment_record_mtc", mr)
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
					"工蜂小智 - 会员套餐",
					"工蜂小智 - 会员套餐",
					convert.ToString(tradeNo),
					convert.MustFloat64(mtc["price"]))
			} else if payType == enum.PAY_TYPE_WECHAT {
				codeUrl, payErr = wxpay.GetPayUrl(
					"工蜂小智 - 会员套餐",
					"工蜂小智 - 会员套餐",
					convert.ToString(tradeNo),
					convert.MustFloat64(mtc["price"]))
			}
			if payErr != nil {
				return false, 0, "", ""
			}
			updateOrderUrl(tradeNo, codeUrl, getUrl)
			return flog, tradeNo, codeUrl, getUrl
		}
		return false, 0, "", ""
	}
}
