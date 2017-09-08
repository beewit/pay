package alipay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/query"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/beekit/utils/uhttp"
	"time"
	"github.com/beewit/pay/global"
	"sort"
	"strings"
)

// Trade trade
type Trade struct{}

// NewTrade new trade
func NewTrade() *Trade {
	return &Trade{}
}

// Sign trade sign
func (t Trade) Sign(args interface{}, privatePath string) (string, error) {
	params, err := query.Values(args)
	if err != nil {
		return "", err
	}
	query, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return "", err
	}
	privateKey := utils.ReadByte(privatePath)

	sign, err := encrypt.NewRsae().RSASign(query, privateKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s&sign=%s",
		query,
		url.QueryEscape(sign),
	), nil
}

// Verify verify
func (t Trade) Verify(args url.Values, publicPath string) error {
	sign := args.Get("sign")
	args.Del("sign")
	args.Del("sign_type")

	var keys = make([]string, 0, 0)
	for key, value := range args {
		if key == "sign" || key == "sign_type" {
			continue
		}
		if len(value) > 0 {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = args.Get(key)
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var s = strings.Join(pList, "&")
	println(s)
	publicKey := utils.ReadByte(publicPath)

	ok, err := encrypt.NewRsae().RSAVerify(s, sign, publicKey)
	if !ok {
		if err != nil {
			return errors.New("签名错误" + err.Error())
		}
		return errors.New("签名错误")
	}
	return nil
}

// Query query
func (t Trade) Query(str string) (Query, error) {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "GET",
		URL:    fmt.Sprintf("https://openapi.alipay.com/gateway.do?%s", str),
	})
	if err != nil {
		return Query{}, err
	}
	result := Query{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.Code != "10000" {
		return result, errors.New(result.Msg)
	}
	return result, nil
}

// Refund refund
func (t Trade) Refund(str string) (Query, error) {
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "GET",
		URL:    fmt.Sprintf("https://openapi.alipay.com/gateway.do?%s", str),
	})
	if err != nil {
		return Query{}, err
	}
	result := Query{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.Code != "10000" {
		return result, errors.New(result.Msg)
	}
	return result, nil
}

func generateTimestampStr() string {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func GetPayUrl(body, subject, tradeNo string, amount float64) (string, string, error) {
	codeUrl, err := sign(Request{
		Body:        body,
		Subject:     subject,
		OutTradeNo:  tradeNo,
		TotalAmount: amount,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		QrPayMode:   "4",
		QrcodeWidth: "200",
		TimeExpress: "1d",
	})
	if err != nil {
		return "", "", err
	}
	getUrl, err2 := sign(Request{
		Body:        body,
		Subject:     subject,
		OutTradeNo:  tradeNo,
		TotalAmount: amount,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		TimeExpress: "1d",
	})
	if err2 != nil {
		return "", "", err
	}
	return codeUrl, getUrl, nil
}

//请求参数qr_pay_mode设置为4 返回二维码
func sign(args Request) (string, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	// 签名
	sign, err := NewTrade().Sign(Sign{
		AppID:      global.AlipayAppId,
		Method:     "alipay.trade.page.pay",
		ReturnURL:  "http://pay.tbqbz.com/app/page/notify/load.html",
		Charset:    "utf-8",
		SignType:   "RSA2",
		TimeStamp:  generateTimestampStr(),
		Version:    "1.0",
		NotifyURL:  "http://pay.tbqbz.com/alipay/notify",
		BizContent: string(body),
	}, global.AlipayPrivatePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?%s",
		"https://openapi.alipay.com/gateway.do",
		sign,
	), nil
}
