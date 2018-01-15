package handler

import (
	"errors"
	"fmt"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/pay/global"
	"github.com/beewit/pay/wxpay"
	"github.com/labstack/echo"
	"strings"
)

func GetRedPacket(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_send_red_packet WHERE status=? AND id=? LIMIT 1", enum.NORMAL, id)
	if err != nil {
		global.Log.Error("GetRedPacket sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

/**
* 发送红包支付
 */
func RedPacketPay(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id == "" {
		return utils.ErrorNull(c, "红包id不能为空")
	}
	openId := strings.TrimSpace(c.FormValue("openId"))
	if openId == "" {
		return utils.ErrorNull(c, "用户openId不能为空")
	}
	if !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "红包id格式错误")
	}
	redPacket := GetRedPacket(convert.MustInt64(id))
	if redPacket == nil {
		return utils.ErrorNull(c, "红包记录不存在")
	}
	totalPrice := convert.MustFloat64(redPacket["money"])
	currentTime := utils.CurrentTime()
	ip := c.RealIP()
	flog := true
	tradeNo := utils.ID()
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		//1、创建支付记录
		//2、创建临时二维码
		order := map[string]interface{}{
			"id":           tradeNo,
			"account_id":   acc.ID,
			"account_name": acc.Nickname,
			"type":         enum.PAY_TYPE_RED_PACKET,
			"pay_type":     enum.PAY_TYPE_WECHAT_MINI_APP,
			"pay_price":    totalPrice,
			"pay_status":   enum.PAY_STATUS_NOT,
			"status":       enum.NORMAL,
			"ct_time":      currentTime,
			"ct_ip":        ip,
		}
		_, err := tx.InsertMap("order_payment", order)
		if err != nil {
			panic(err)
		}
		orderRecord := map[string]interface{}{
			"id":                         utils.ID(),
			"order_payment_id":           tradeNo,
			"account_send_red_packet_id": id,
			"price": totalPrice,
		}
		_, err = tx.InsertMap("order_payment_record_red_packet", orderRecord)
		if err != nil {
			panic(err)
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("创建红包支付订单失败，%v", err)
			flog = false
		}
	})
	if err != nil {
		return utils.ErrorNull(c, "创建红包支付订单失败")
	} else {
		//支付接口
		body := "工蜂引流 - 发红包"
		subject := "工蜂引流 - 发红包"
		defray, err := wxpay.GetMiniAppPayPars(body, subject, convert.ToString(tradeNo), openId, totalPrice)
		if err != nil {
			return utils.Error(c, "创建支付签名失败:"+err.Error(), nil)
		}
		return utils.Success(c, "创建红包支付订单成功", map[string]interface{}{
			"tradeNo":    tradeNo,
			"totalPrice": totalPrice,
			"sign":       defray.Sign,
			"appId":      global.WechatMiniAppConf.AppID,
			"partnerId":  global.WechatMchID,
			"prepayId":   defray.PrepayID,
			"noncestr":   defray.NonceStr,
			"timeStamp":  defray.TimeStamp})
	}
}

//红包支付成功回调
func UpdateOrderRedPacketStatus(order map[string]interface{}, price float64, ip string) bool {
	id := convert.MustInt64(order["id"])
	var x int64
	var err error
	var errMsg string
	flog := true
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		sql := "UPDATE order_payment SET pay_status=?,pay_time=?,pay_ip=? WHERE id=?"
		x, err = tx.Update(sql, enum.PAY_STATUS_END, utils.CurrentTime(), ip, id)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		if x <= 0 {
			errMsg = fmt.Sprintf("%v修改订单状态失败", id)
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}
		recordList := GetOrderRecordRedPacketList(id)
		if recordList == nil {
			errMsg = fmt.Sprintf("%v获取订单红包记录失败", id)
			global.Log.Error(errMsg)
			panic(errors.New(errMsg))
		}
		var qrCodePath string
		var redPacketId int64
		for i := 0; i < len(recordList); i++ {
			qrCodePath = ""
			redPacketId = convert.MustInt64(recordList[0]["id"])
			body, err := uhttp.Cmd(uhttp.Request{
				Method: "GET",
				URL:    fmt.Sprintf("http://m.9ee3.com/account/create/temporary/qrcode?objId=%v&objType=%s", redPacketId, enum.QRCODE_RED_PACKET),
			})
			if err != nil {
				global.Log.Error("获取领取红包临时二维码失败，%v", err.Error())
			} else {
				//global.Log.Info(string(body))
				resultParam := utils.ToResultParam(body)
				if resultParam.Ret == utils.SUCCESS_CODE {
					data, err := convert.Obj2Map(resultParam.Data)
					if err != nil {
						global.Log.Error("获取领取红包临时二维码失败，转换数据失败：%v", err.Error())
					} else {
						//保存
						qrCodePath = convert.ToString(data["path"])
					}
				} else {
					global.Log.Error("获取领取红包临时二维码失败，%v", resultParam.Msg)
				}
			}
			x, err = tx.Update("UPDATE account_send_red_packet SET qrcode=?,pay_state=? WHERE id=?", qrCodePath, enum.PAY_STATUS_END, redPacketId)
			if err != nil {
				global.Log.Error(err.Error())
				panic(err)
			}
			if x <= 0 {
				errMsg = fmt.Sprintf("%v修改红包支付状态失败", redPacketId)
				global.Log.Error(errMsg)
				panic(errors.New(errMsg))
			}
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("支付订单通知处理失败，%v", err)
			flog = false
		}
	})
	return flog
}

func GetOrderRecordRedPacketList(id int64) []map[string]interface{} {
	maps, err := global.DB.Query("SELECT red.* FROM order_payment_record_red_packet rec LEFT JOIN account_send_red_packet red ON rec.account_send_red_packet_id=red.id WHERE order_payment_id=?", id)
	if err != nil {
		global.Log.Error("GetOrderRecordRedPacket sql error:%s", err.Error())
		return nil
	}
	return maps
}
