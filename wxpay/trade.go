package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/beewit/beekit/utils/query"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/beekit/utils/uhttp"
	"strconv"
	"time"
	"crypto/md5"
	"github.com/beewit/pay/global"
	"github.com/beewit/beekit/utils"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"image/png"
	"bytes"
	"github.com/beewit/beekit/utils/imgbase64"

)

// Trade trade
type Trade struct{}

// NewTrade new trade
func NewTrade() *Trade {
	return &Trade{}
}

// Sign trade sign
func (t Trade) Sign(args interface{}, key string) (string, error) {
	params, err := query.Values(args)
	if err != nil {
		return "", err
	}
	query, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return "", err
	}
	global.Log.Info(query)
	sign := encrypt.NewRsae().Md532(
		fmt.Sprintf("%s&key=%s",
			query,
			key,
		),
	)
	return strings.ToUpper(sign), nil
}

// Prepay trade perpay
func (t Trade) Prepay(args Sign) (Prepay, error) {
	body, err := xml.Marshal(args)
	if err != nil {
		return Prepay{}, err
	}
	println(string(body[:]))
	header := http.Header{}
	header.Add("Accept", "application/xml")
	header.Add("Content-Type", "application/xml;charset=utf-8")
	body, err = uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    "https://api.mch.weixin.qq.com/pay/unifiedorder",
		Body:   body,
		Header: header,
	})
	if err != nil {
		return Prepay{}, err
	}
	result := Prepay{}
	println(string(body[:]))
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.ReturnCode != Success {
		return result, errors.New(result.ReturnMsg)
	}
	if result.ResultCode != Success {
		return result, errors.New(result.ErrCodeDes)
	}
	return result, nil
}

// Verify verify
func (t Trade) Verify(args Notice, key string) error {
	//args, err := query.Values(n)
	//if err != nil {
	//	return err
	//}
	//rSign := args.Get("sign")
	//args.Del("sign")
	//
	//var keys = make([]string, 0, 0)
	//for key, value := range args {
	//	if len(value) > 0 {
	//		keys = append(keys, key)
	//	}
	//}
	//
	//sort.Strings(keys)
	//
	//var pList = make(map[string]string)
	//for _, key := range keys {
	//	var value = args.Get(key)
	//	if len(value) > 0 {
	//		pList[key] = value
	//	}
	//}

	sign, err := t.Sign(args, key)
	if err != nil {
		return err
	}
	if args.Sign != sign {
		return errors.New("签名错误")
	}
	return nil
}

// Query query
func (t Trade) Query(args Sign) (Query, error) {
	body, err := xml.Marshal(args)
	if err != nil {
		return Query{}, err
	}
	header := http.Header{}
	header.Add("Accept", "application/xml")
	header.Add("Content-Type", "application/xml;charset=utf-8")
	body, err = uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    "https://api.mch.weixin.qq.com/pay/orderquery",
		Body:   body,
		Header: header,
	})
	if err != nil {
		return Query{}, err
	}
	result := Query{}
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.ReturnCode != Success {
		return result, errors.New(result.ReturnMsg)
	}
	if result.ResultCode != Success {
		return result, errors.New(result.ErrCodeDes)
	}
	return result, nil
}

// Refund refund
func (t Trade) Refund(args Sign) (Refund, error) {
	body, err := xml.Marshal(args)
	if err != nil {
		return Refund{}, err
	}
	header := http.Header{}
	header.Add("Accept", "application/xml")
	header.Add("Content-Type", "application/xml;charset=utf-8")
	body, err = uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    "https://api.mch.weixin.qq.com/pay/refund",
		Body:   body,
		Header: header,
	})
	if err != nil {
		return Refund{}, err
	}
	result := Refund{}
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	if result.ReturnCode != Success {
		return result, errors.New(result.ReturnMsg)
	}
	if result.ResultCode != Success {
		return result, errors.New(result.ErrCodeDes)
	}
	return result, nil
}

func GenerateNonceStr() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

func GetCodeUrl(args Request) (string, error) {
	sign := Sign{
		AppID:          global.WechatAppId,
		MchID:          global.WechatMchID,
		NonceStr:       GenerateNonceStr(),
		TradeType:      "NATIVE",
		SpbillCreateIP: utils.GetIp(),
		NotifyURL:      global.WechatNotifyURL,
		Request:        args,
	}
	str, err := NewTrade().Sign(sign, global.WechatApiKey)
	if err != nil {
		return "", err
	}
	sign.Sign = str
	prepay, err := NewTrade().Prepay(sign)
	if err != nil {
		return "", err
	}
	return prepay.CodeURL, nil
}

func GetPayUrl(body, subject, tradeNo string, amount float64) (string, error) {
	r := Request{
		Body:       body,
		Attach:     subject,
		OutTradeNo: tradeNo,
		ProductID:  tradeNo,
		TotalFee:   (int)(amount * 100),
	}
	codeUrl, err := GetCodeUrl(r)
	if err != nil {
		global.Log.Error("微信扫码支付GetPayUrl：", err.Error())
		return "", err
	}
	return CreateQrCode(codeUrl)
}

func CreateQrCode(url string) (string, error) {
	global.Log.Info(url)
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		global.Log.Error("微信扫码支付CreateQrCode：", err.Error())
		return "", err
	}
	// Scale the barcode to 300*300 pixelscustomer_educational
	qrCode, err = barcode.Scale(qrCode, 300, 300)
	if err != nil {
		global.Log.Error("微信扫码支付CreateQrCode：", err.Error())
		return "", err
	}
	var b bytes.Buffer
	png.Encode(&b, qrCode)

	return imgbase64.FromBuffer(b), nil
}
