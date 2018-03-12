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

//多个功能开通
func CreateBatchFuncOrder(c echo.Context) error {
	funcIdStr := c.FormValue("funcId")
	fcIdStr := c.FormValue("fcId")
	pt := c.FormValue("pt")
	payWalletMoneyStr := c.FormValue("walletMoney")
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
	var payWalletMoney float64 = 0
	if payWalletMoneyStr != "" && utils.IsValidNumber(payWalletMoneyStr) {
		payWalletMoney = convert.MustFloat64(payWalletMoneyStr)
		m := getWallet(acc.ID)
		if m == nil {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
		walletMoney := convert.MustFloat64(m["money"])
		if payWalletMoney > walletMoney {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
	}
	accID := acc.ID
	accName := acc.Nickname
	flog, tradeNo, _, codeUrl, getUrl := getOrderCode(mt, fc, accID, accName, pt, c.RealIP(), payWalletMoney, true)
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
	payWalletMoneyStr := c.FormValue("walletMoney")
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
		return utils.Error(c, "选择的功能开通方式不存在", nil)
	}
	acc := global.ToInterfaceAccount(c.Get("account"))
	var payWalletMoney float64 = 0
	if payWalletMoneyStr != "" && utils.IsValidNumber(payWalletMoneyStr) {
		payWalletMoney = convert.MustFloat64(payWalletMoneyStr)
		m := getWallet(acc.ID)
		if m == nil {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
		walletMoney := convert.MustFloat64(m["money"])
		if payWalletMoney > walletMoney {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
	}
	accID := acc.ID
	accName := acc.Nickname
	flog, tradeNo, _, codeUrl, getUrl := getOrderCode(mt, fc, accID, accName, pt, c.RealIP(), payWalletMoney, true)
	if flog {
		return utils.Success(c, "创建订单成功", map[string]interface{}{"codeUrl": codeUrl, "getUrl": getUrl, "tradeNo": tradeNo})
	}
	return utils.Error(c, "创建订单失败", nil)
}

//只是创建订单、不创建支付二维码
func CreateAppOrder(c echo.Context) error {
	body := c.FormValue("body")
	subject := c.FormValue("subject")
	funcIdStr := c.FormValue("funcId")
	fcIdStr := c.FormValue("fcId")
	pt := c.FormValue("pt")
	payWalletMoneyStr := c.FormValue("walletMoney")
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
	var payWalletMoney float64 = 0
	if payWalletMoneyStr != "" && utils.IsValidNumber(payWalletMoneyStr) {
		payWalletMoney = convert.MustFloat64(payWalletMoneyStr)
		m := getWallet(acc.ID)
		if m == nil {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
		walletMoney := convert.MustFloat64(m["money"])
		if payWalletMoney > walletMoney {
			return utils.ErrorNull(c, "余额支付的金额已超过可用的钱包余额")
		}
	}
	accID := acc.ID
	accName := acc.Nickname
	flog, tradeNo, totalPrice, _, _ := getOrderCode(mt, fc, accID, accName, pt, c.RealIP(), payWalletMoney, false)
	if flog {
		if pt == enum.PAY_TYPE_WECHATAPP {
			defray, err := wxpay.GetAppPayPars(body, subject, convert.ToString(tradeNo), totalPrice)
			if err != nil {
				return utils.Error(c, "创建支付签名失败", nil)
			}
			return utils.Success(c, "创建订单成功", map[string]interface{}{
				"tradeNo":    tradeNo,
				"totalPrice": totalPrice,
				"sign":       defray.Sign,
				"appId":      global.WechatAPPAppId,
				"partnerId":  global.WechatAPPMchID,
				"prepayId":   defray.PrepayID,
				"noncestr":   defray.NonceStr,
				"timeStamp":  defray.TimeStamp})
		} else if pt == enum.PAY_TYPE_ALIPAY {
			sign, err := alipay.GetAPPSign(alipay.Request{
				Body:        body,
				Subject:     subject,
				OutTradeNo:  convert.ToString(tradeNo),
				TotalAmount: totalPrice,
				ProductCode: "QUICK_MSECURITY_PAY",
				TimeExpress: "1d",
			})
			if err != nil {
				return utils.Error(c, "创建支付签名失败", nil)
			}
			return utils.Success(c, "创建订单成功", map[string]interface{}{"tradeNo": tradeNo, "totalPrice": totalPrice, "sign": sign})
		} else if pt == enum.PAY_TYPE_WECHAT_MINI_APP {
			ws, err := GetMiniAppSession(c)
			if err != nil {
				return utils.AuthFailNull(c)
			}
			defray, err := wxpay.GetMiniAppPayPars(body, subject, convert.ToString(tradeNo), ws.Openid, totalPrice)
			if err != nil {
				return utils.Error(c, "创建支付签名失败:"+err.Error(), nil)
			}
			return utils.Success(c, "创建支付订单成功", map[string]interface{}{
				"sign":      defray.Sign,
				"package":   defray.Package,
				"noncestr":  defray.NonceStr,
				"timeStamp": defray.TimeStamp,
				"tradeNo":   tradeNo,
			})
		} else if pt == enum.PAY_TYPE_WECHAT_H5 {
			u := GetOauthUser(c)
			if u == nil {
				return utils.AuthWechatFailNull(c)
			}
			defray, err := wxpay.GetMPPayPars(body, subject, convert.ToString(tradeNo), u.OpenId, totalPrice)
			if err != nil {
				return utils.Error(c, "创建支付签名失败:"+err.Error(), nil)
			}
			return utils.Success(c, "创建红包支付订单成功", map[string]interface{}{
				"sign":      defray.PaySign,
				"package":   defray.Package,
				"noncestr":  defray.NonceStr,
				"timeStamp": defray.TimeStamp,
				"tradeNo":   tradeNo,
			})
		} else {
			return utils.Error(c, "当前仅支持微信和支付宝支付", nil)
		}
	}
	return utils.Error(c, "创建订单失败", nil)
}

func getWallet(accId int64) map[string]interface{} {
	sql := "SELECT * FROM account_wallet WHERE account_id=? LIMIT 1"
	rows, _ := global.DB.Query(sql, accId)
	if rows == nil || len(rows) != 1 {
		return nil
	}
	return rows[0]
}

//订单继续支付
func OrderPay(c echo.Context) error {
	acc := global.ToInterfaceAccount(c.Get("account"))
	orderId := c.FormValue("orderId")
	body := c.FormValue("body")
	subject := c.FormValue("subject")
	if orderId == "" || !utils.IsValidNumber(orderId) {
		return utils.ErrorNull(c, "订单号无效")
	}
	order := GetOrder(convert.MustInt64(orderId))
	if order == nil {
		return utils.ErrorNull(c, "订单不存在")
	}
	if acc.ID != convert.MustInt64(order["account_id"]) {
		return utils.ErrorNull(c, "无订单操作权限")
	}
	//订单正常并且未支付状态才可取消订单
	if convert.ToString(order["status"]) != enum.NORMAL && convert.ToString(order["pay_status"]) != enum.PAY_STATUS_NOT {
		return utils.ErrorNull(c, "订单非有效状态无法继续支付")
	}
	payPrice := convert.MustFloat64(order["pay_price"])
	pt := convert.ToString(order["pay_type"])
	if pt == enum.PAY_TYPE_WECHATAPP {
		defray, err := wxpay.GetAppPayPars(body, subject, orderId, payPrice)
		if err != nil {
			return utils.Error(c, "创建支付签名失败", nil)
		}
		return utils.Success(c, "继续订单支付成功", map[string]interface{}{
			"tradeNo":    orderId,
			"totalPrice": payPrice,
			"sign":       defray.Sign,
			"appId":      global.WechatAPPAppId,
			"partnerId":  global.WechatAPPMchID,
			"prepayId":   defray.PrepayID,
			"noncestr":   defray.NonceStr,
			"timeStamp":  defray.TimeStamp})
	} else if pt == enum.PAY_TYPE_ALIPAY {
		sign, err := alipay.GetAPPSign(alipay.Request{
			Body:        body,
			Subject:     subject,
			OutTradeNo:  orderId,
			TotalAmount: payPrice,
			ProductCode: "QUICK_MSECURITY_PAY",
			TimeExpress: "1d",
		})
		if err != nil {
			return utils.Error(c, "创建支付签名失败", nil)
		}
		return utils.Success(c, "继续订单支付成功", map[string]interface{}{"tradeNo": orderId, "totalPrice": payPrice, "sign": sign})
	} else if pt == enum.PAY_TYPE_WECHAT_MINI_APP {
		ws, err := GetMiniAppSession(c)
		if err != nil {
			return utils.AuthFailNull(c)
		}
		defray, err := wxpay.GetMiniAppPayPars(body, subject, orderId, ws.Openid, payPrice)
		if err != nil {
			return utils.Error(c, "创建支付签名失败:"+err.Error(), nil)
		}
		return utils.Success(c, "继续订单支付成功", map[string]interface{}{
			"sign":      defray.Sign,
			"package":   defray.Package,
			"noncestr":  defray.NonceStr,
			"timeStamp": defray.TimeStamp,
			"tradeNo":   orderId,
		})
	} else if pt == enum.PAY_TYPE_WECHAT_H5 {
		u := GetOauthUser(c)
		if u == nil {
			return utils.AuthWechatFailNull(c)
		}
		defray, err := wxpay.GetMPPayPars(body, subject, orderId, u.OpenId, payPrice)
		if err != nil {
			return utils.Error(c, "创建支付签名失败:"+err.Error(), nil)
		}
		return utils.Success(c, "创建红包支付订单成功", map[string]interface{}{
			"sign":      defray.PaySign,
			"package":   defray.Package,
			"noncestr":  defray.NonceStr,
			"timeStamp": defray.TimeStamp,
			"tradeNo":   orderId,
		})
	} else {
		return utils.Error(c, "当前仅支持微信和支付宝支付", nil)
	}
}

func CancelOrder(c echo.Context) error {
	acc := global.ToInterfaceAccount(c.Get("account"))
	orderId := c.FormValue("orderId")
	if orderId == "" || !utils.IsValidNumber(orderId) {
		return utils.ErrorNull(c, "订单号无效")
	}
	order := GetOrder(convert.MustInt64(orderId))
	if order == nil {
		return utils.ErrorNull(c, "订单不存在")
	}
	if acc.ID != convert.MustInt64(order["account_id"]) {
		return utils.ErrorNull(c, "无订单操作权限")
	}
	//订单正常并且未支付状态才可取消订单
	if convert.ToString(order["status"]) != enum.NORMAL && convert.ToString(order["pay_status"]) != enum.PAY_STATUS_NOT {
		return utils.ErrorNull(c, "订单不是有效状态或订单不是非支付状态无法取消")
	}
	flog := false
	var err error
	payMoney := convert.MustFloat64(order["pay_money"])
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		//余额支付判断
		if payMoney > 0 {
			//需要退换余额
			wallet := getWallet(acc.ID)
			if wallet == nil {
				err = errors.New("账号余额表无记录")
				panic(err)
			} else {
				changeMoney := convert.MustFloat64(wallet["money"]) + payMoney
				_, err = tx.Update("UPDATE account_wallet SET money=money+? WHERE account_id=?", payMoney, acc.ID)
				if err != nil {
					global.Log.Error(err.Error())
					panic(err)
				}
				//添加余额日志记录
				_, err = tx.InsertMap("account_wallet_log", map[string]interface{}{
					"id":               utils.ID(),
					"account_id":       acc.ID,
					"change_money":     payMoney,
					"money":            changeMoney,
					"order_payment_id": orderId,
					"type":             enum.WALLET_CANCEL_PAY,
					"remark":           "取消订单返回订单的余额支付金额",
					"ct_time":          utils.CurrentTime(),
				})
				if err != nil {
					global.Log.Error(err.Error())
					panic(err)
				}
			}
		}
		_, err = tx.Update("UPDATE order_payment SET status=? WHERE id=?", enum.CANCEL, orderId)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		flog = true
	}, func(err error) {
		if err != nil {
			global.Log.Error("取消订单失败，ERROR：%v", err)
			flog = false
		}
	})
	if flog {
		if payMoney > 0 {
			return utils.SuccessNull(c, fmt.Sprintf("取消订单成功并返回余额 ¥%v", payMoney))
		}
		return utils.SuccessNull(c, "取消订单成功")
	}
	return utils.ErrorNull(c, "取消订单失败")
}

func GetFuncAndCharge(c echo.Context) error {
	fid := c.FormValue("fid")
	if fid == "" {
		return utils.ErrorNull(c, "请选择开通功能")
	}
	t := c.FormValue("type")
	switch t {
	case convert.ToString(enum.FUNC_CHARGE_1):
	case convert.ToString(enum.FUNC_CHARGE_2):
		break
	default:
		t = "0"
		break
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
	sql = "SELECT * FROM func_charge WHERE status=? AND type=? ORDER BY `order` DESC,ct_time DESC"
	rows2, err2 := global.DB.Query(sql, enum.NORMAL, t)
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
	id := convert.MustInt64(c.FormValue("out_trade_no"))
	totalMoney := convert.MustFloat64(c.FormValue("total_amount"))
	ip := c.RealIP()
	order := GetOrder(id)
	UpdateOrderStatus(order, totalMoney, ip)
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
	id := convert.MustInt64(args.OutTradeNo)
	if id <= 0 {
		global.Log.Error("未获取到有效订单")
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  "未获取到有效订单",
		}
		return c.XML(http.StatusOK, wr)
	}
	order := GetOrder(id)
	if order == nil {
		global.Log.Error("未查询到该订单:%v", id)
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  fmt.Sprintf("未查询到该订单:%v", id),
		}
		return c.XML(http.StatusOK, wr)
	}
	var apiKey string
	if order["pay_type"] == enum.PAY_TYPE_WECHATAPP {
		apiKey = global.WechatAPPApiKey
	} else if order["pay_type"] == enum.PAY_TYPE_WECHAT || order["pay_type"] == enum.PAY_TYPE_WECHAT_MINI_APP || order["pay_type"] == enum.PAY_TYPE_WECHAT_H5 {
		apiKey = global.WechatApiKey
	}
	err := wxpay.NewTrade().Verify(args, apiKey)
	if err != nil {
		global.Log.Error(err.Error())
		wr := &wxpay.Response{
			ReturnCode: "ERROR",
			ReturnMsg:  "验签失败",
		}
		return c.XML(http.StatusOK, wr)
	}
	totalMoney := convert.MustFloat64(args.TotalFee) / 100
	ip := c.RealIP()
	UpdateOrderStatus(order, totalMoney, ip)
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
	order := GetOrder(convert.MustInt64(tradeNo))
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

func UpdateOrderStatus(order map[string]interface{}, price float64, ip string) bool {
	id := convert.MustInt64(order["id"])
	var errMsg string
	ids := convert.ToString(id)
	if order == nil {
		errMsg = "订单不存在：" + ids
		global.Log.Error(errMsg)
		return false
	}
	if order["pay_status"] == enum.PAY_STATUS_END {
		errMsg = "订单已支付过了：" + ids
		global.Log.Error(errMsg)
		return false
	}
	if convert.MustFloat64(order["pay_price"]) != price {
		errMsg = "订单支付金额不一致"
		global.Log.Warning(errMsg)
		return false
	}
	global.Log.Info("订单类型：%v", order["type"])
	switch order["type"] {
	//功能开通
	case enum.PAY_TYPE_FUNC:
		return UpdateOrderFuncStatus(order, price, ip)
	case enum.PAY_TYPE_RED_PACKET:
		return UpdateOrderRedPacketStatus(order, price, ip)
	}
	global.Log.Error("无此订单开通类型：%v", order["type"])
	return false
}

func UpdateOrderFuncStatus(order map[string]interface{}, price float64, ip string) bool {
	var x int64
	var err error
	var errMsg string
	flog := false
	id := convert.MustInt64(order["id"])
	accId, errAccId := convert.ToInt64(order["account_id"])
	if errAccId != nil {
		global.Log.Error("%v订单付款的帐号异常：%s", id, errAccId.Error())
		return false
	}
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		sql := "UPDATE order_payment SET pay_status=?,pay_time=?,pay_ip=? WHERE id=?"
		x, err = tx.Update(sql, enum.PAY_STATUS_END, utils.CurrentTime(), ip, id)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			errMsg := fmt.Sprintf("%v修改订单状态失败：", id)
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}

		var daysTime time.Time
		//订单开通功能
		orderFunc := getOrderFuncId(id)
		//帐号已开通功能记录
		accFunc := getAccountFuncByAccId(accId)
		acMaps := []map[string]interface{}{}
		for i := 0; i < len(orderFunc); i++ {

			//测试数据只能为指定账户才能使用
			var testFuncChargeId int64 = 7
			if convert.MustInt64(orderFunc[i]["func_charge_id"]) == testFuncChargeId && accId != 122068319091036160 {
				global.Log.Error("非指定账户不能使用测试功能支付，支付无效")
				panic(errMsg)
			}

			flog := false
			days := convert.MustInt(orderFunc[i]["days"])
			giveDays := convert.MustInt(orderFunc[i]["give_days"])
			if giveDays > 0 {
				days += giveDays
			}
			funcId := orderFunc[i]["func_id"]

			for j := 0; j < len(accFunc); j++ {
				if accFunc[j]["func_id"] == orderFunc[i]["func_id"] {
					if accFunc[j]["expiration_time"] != nil {
						expirTimeStr, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(accFunc[j]["expiration_time"]))
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
					} else {
						//无到期时间
						daysTime = time.Now().AddDate(0, 0, days)
					}
					flog = true
					break
				}
			}
			if flog {
				//修改
				sql = "UPDATE account_func SET expiration_time=?,ut_time=? WHERE account_id=? AND func_id=?"
				x, err = tx.Update(sql, utils.FormatTime(daysTime), utils.CurrentTime(), accId, funcId)
				if err != nil {
					global.Log.Error(err.Error())
					panic(err)
				}
			} else {
				//添加
				daysTime = time.Now().AddDate(0, 0, days)
				m := make(map[string]interface{})
				m["id"] = utils.ID()
				m["account_id"] = accId
				m["func_id"] = orderFunc[i]["func_id"]
				m["expiration_time"] = utils.FormatTime(daysTime)
				m["ct_time"] = utils.CurrentTime()
				m["ut_time"] = utils.CurrentTime()
				acMaps = append(acMaps, m)
			}
		}
		if len(acMaps) > 0 {
			x, err = tx.InsertMapList("account_func", acMaps)
			if err != nil {
				global.Log.Error(err.Error())
				panic(err)
			}
		}

		//查询当前用户是否有邀请者，如果有则进行支付金额按比例进行返利，四舍五入
		orderAccount := getAccountId(accId)
		if orderAccount != nil && orderAccount["share_account_id"] != nil && convert.MustInt64(orderAccount["share_account_id"]) > 0 {
			//对分享者加入钱包记录和钱包增加 share_account_id
			shareAccountId := convert.MustInt64(orderAccount["share_account_id"])
			shareAccount := getAccountId(shareAccountId)
			//判断是否从该用户出获得过邀请返利
			shareRebateWalletLog := getShareRebateWalletLog(shareAccountId, accId, enum.WALLET_REBATE)
			if shareRebateWalletLog != nil {
				global.Log.Warning("已经获得过邀请初次返利了，暂不重复下单返利")
			}
			if shareAccount != nil && shareRebateWalletLog == nil {
				//邀请返利，暂定百分之五
				changeMoney := convert.MustFloat64(fmt.Sprintf("%.2f", price*0.05))
				accountWallet := getAccountWalletByAccId(shareAccountId)
				var money = changeMoney
				if accountWallet != nil {
					money = convert.MustFloat64(accountWallet["money"]) + changeMoney
				}
				wallet := getAccountWalletByAccId(accId)
				if wallet == nil {
					sql := "INSERT INTO account_wallet(id,account_id,money,last_time,last_ip)VALUES(?,?,?,?,?)"
					_, err = tx.Insert(sql, utils.ID(), shareAccountId, money, utils.CurrentTime(), ip)
					if err != nil {
						global.Log.Error(fmt.Sprintf("订单通知，邀请者获得返利时修改钱包金额失败，错误：%s", err.Error()))
						panic(err)
					}
				} else {
					sql := "UPDATE account_wallet SET money=?,last_time=?,last_ip=? WHERE account_id=?"
					_, err = tx.Update(sql, money, utils.CurrentTime(), ip, shareAccountId)
					if err != nil {
						global.Log.Error(fmt.Sprintf("订单通知，邀请者获得返利时修改钱包金额失败，错误：%s", err.Error()))
						panic(err)
					}
				}
				_, err = tx.InsertMap("account_wallet_log", map[string]interface{}{
					"id":                   utils.ID(),
					"account_id":           shareAccountId,
					"change_money":         changeMoney,
					"money":                money,
					"order_payment_id":     id,
					"type":                 enum.WALLET_REBATE,
					"was_share_account_id": accId,
					"remark":               fmt.Sprintf("邀请会员[%v]获得返利专享！", accId),
					"ct_time":              utils.CurrentTime(),
					"ct_ip":                ip,
				})
				if err != nil {
					global.Log.Error(fmt.Sprintf("订单通知，邀请者获得返利时添加钱包日志失败，错误：%s", err.Error()))
					panic(err)
				}
			}
		}

		flog = true
		errMsg = fmt.Sprintf("%v修改订单成功", id)
		global.Log.Info(errMsg)
	}, func(err error) {
		if err != nil {
			global.Log.Error("支付订单通知处理失败，%v", err)
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

//该帐号开通的功能项目
func getAccountFuncByAccId(accId int64) []map[string]interface{} {
	sql := "SELECT * FROM account_func WHERE account_id=?"
	rows, err := global.DB.Query(sql, accId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows
}

func getAccountWalletByAccId(accId int64) map[string]interface{} {
	sql := "SELECT * FROM account_wallet WHERE account_id=?"
	rows, err := global.DB.Query(sql, accId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

//邀请返利记录
func getShareRebateWalletLog(accId, wasShareAccId int64, t string) map[string]interface{} {
	sql := "SELECT * FROM account_wallet_log WHERE account_id=? AND was_share_account_id=? AND type=?"
	rows, err := global.DB.Query(sql, accId, wasShareAccId, t)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

func getOrderFuncId(orderId int64) []map[string]interface{} {
	sql := "SELECT * FROM order_payment_record_func WHERE order_payment_id=?"
	rows, err := global.DB.Query(sql, orderId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows
}

func GetOrder(id int64) map[string]interface{} {
	sql := "SELECT * FROM order_payment o WHERE o.id=? LIMIT 1"
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

func getOrderCode(mt []map[string]interface{}, fc map[string]interface{}, accId int64, accName, payType, ip string, payWalletMoney float64, isPayQrCode bool) (bool, int64, float64, string,
	string) {
	m := make(map[string]interface{})
	tradeNo := utils.ID()
	var sumPrice float64 = 0
	funcList := make([]map[string]interface{}, len(mt))
	for i := 0; i < len(mt); i++ {
		mr := make(map[string]interface{})
		mr["id"] = utils.ID()
		mr["order_payment_id"] = tradeNo
		mr["func_id"] = mt[i]["id"]
		mr["func_name"] = mt[i]["name"]
		mr["price"] = mt[i]["price"]
		mr["func_charge_id"] = fc["id"]
		mr["days"] = fc["days"]
		mr["give_days"] = fc["give_days"]
		funcList[i] = mr

		sumPrice += convert.MustFloat64(mr["price"])
	}

	discount := convert.MustFloat64(fc["discount"])

	totalPrice := convert.MustFloat64(fc["days"]) * sumPrice
	if discount > 0 {
		totalPrice = totalPrice * discount

	}
	totalPrice = convert.MustFloat64(fmt.Sprintf("%.2f", totalPrice))
	if payWalletMoney > 0 {
		totalPrice = totalPrice - payWalletMoney
	}

	//测试使用
	var testFuncChargeId int64 = 7
	if convert.MustInt64(fc["id"]) == testFuncChargeId {
		totalPrice = 0.01
	}

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
	m["pay_money"] = payWalletMoney

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
		if payWalletMoney > 0 {
			wallet := getWallet(accId)
			if wallet == nil {
				err = errors.New("账号余额表无记录")
				panic(err)
			} else {
				changeMoney := convert.MustFloat64(wallet["money"]) - payWalletMoney
				if changeMoney < 0 {
					//余额不足
					err = errors.New("账号余额不足")
					panic(err)
				} else {
					//扣除余额
					_, err = tx.Update("UPDATE account_wallet SET money=money-? WHERE account_id=?", payWalletMoney, accId)
					if err != nil {
						global.Log.Error(err.Error())
						panic(err)
					}
					//添加余额日志记录
					_, err = tx.InsertMap("account_wallet_log", map[string]interface{}{
						"id":               utils.ID(),
						"account_id":       accId,
						"change_money":     -payWalletMoney,
						"money":            changeMoney,
						"order_payment_id": tradeNo,
						"type":             enum.WALLET_PAY,
						"remark":           "余额订单支付",
						"ct_time":          utils.CurrentTime(),
					})
					if err != nil {
						global.Log.Error(err.Error())
						panic(err)
					}
				}
			}
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
			flog = false
		}
	})

	if flog {
		if isPayQrCode {
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
				return false, 0, 0, "", ""
			}
			updateOrderUrl(tradeNo, codeUrl, getUrl)
			return flog, tradeNo, totalPrice, codeUrl, getUrl
		} else {
			return flog, tradeNo, totalPrice, "", ""
		}
	}
	return false, 0, 0, "", ""
}
