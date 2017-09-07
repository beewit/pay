package main

import (
	"fmt"
	"github.com/beewit/beekit/utils/uhttp"
	"time"
	"testing"
	"github.com/beewit/pay/alipay"
	"encoding/json"
)

func TestPay(t *testing.T) {
	sign, err := sign(alipay.Request{
		Body:        "工蜂小智 - 大数据专家",
		Subject:     "工蜂小智",
		OutTradeNo:  "25889548555442",
		TotalAmount: 0.04,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		QrPayMode:   "4", //扫码二维码
	})
	fmt.Println(err)
	fmt.Println("url：" + sign)

	body, err := uhttp.Cmd(uhttp.Request{
		Method: "GET",
		URL:    sign,
	})
	if err != nil {
		println(err.Error())
	}
	println("成功：" + string(body[:]))
}

//请求参数qr_pay_mode设置为4 返回二维码
func sign(args alipay.Request) (string, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	// 签名
	sign, err := alipay.NewTrade().Sign(alipay.Sign{
		AppID:      "2017082308338334",
		Method:     "alipay.trade.page.pay",
		ReturnURL:  "http://pay.tbqbz.com/alipay/return",
		Charset:    "utf-8",
		SignType:   "RSA2",
		TimeStamp:  generateTimestampStr(),
		Version:    "1.0",
		NotifyURL:  "http://pay.tbqbz.com/alipay/notify",
		BizContent: string(body),
	}, "rsa_private_key.pem")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?%s",
		"https://openapi.alipay.com/gateway.do",
		sign,
	), nil
}

func generateTimestampStr() string {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}
