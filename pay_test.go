package main

import (
	"fmt"
	"github.com/beewit/beekit/utils/uhttp"
	"time"
	"testing"
	"github.com/beewit/pay/alipay"
	"encoding/json"
	"github.com/beewit/pay/wxpay"
	"strconv"
	"github.com/beewit/beekit/utils"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode"
	"os"
	"image/png"
	"github.com/beewit/pay/global"
	"strings"
	"bytes"
	"github.com/beewit/beekit/utils/imgbase64"
	"github.com/beewit/beekit/utils/convert"
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
		AppID:      global.AlipayAppId,
		Method:     "alipay.trade.page.pay",
		ReturnURL:  global.AlipayReturnURL,
		Charset:    "utf-8",
		SignType:   "RSA2",
		TimeStamp:  generateTimestampStr(),
		Version:    "1.0",
		NotifyURL:  global.AlipayNotifyURL,
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

func TestWeiXinPay(t *testing.T) {
	r := wxpay.Request{
		Body:       "工蜂小智",
		Attach:     "工蜂小智慧",
		OutTradeNo: "123456789564555",
		ProductID:  "123456789564555",
		TotalFee:   1,
	}
	wd, err := Sign(r)
	if err != nil {
		println("error：", err.Error())
	}
	println("PrepayID：", wd.PrepayID)
	println("CodeUrl：", wd.CodeUrl)

}

func Sign(args wxpay.Request) (wxpay.Defray, error) {
	sign := wxpay.Sign{
		AppID:          global.WechatAppId,
		MchID:          global.WechatMchID,
		NonceStr:       wxpay.GenerateNonceStr(),
		TradeType:      "NATIVE",
		SpbillCreateIP: utils.GetIp(),
		NotifyURL:      global.WechatNotifyURL,
		Request:        args,
	}
	str, err := wxpay.NewTrade().Sign(sign, global.WechatApiKey)
	if err != nil {
		return wxpay.Defray{}, err
	}
	sign.Sign = str
	prepay, err := wxpay.NewTrade().Prepay(sign)
	if err != nil {
		return wxpay.Defray{}, err
	}
	defray := wxpay.Defray{
		AppID:     global.WechatAppId,
		Package:   "prepay_id=" + prepay.PrepayID,
		NonceStr:  wxpay.GenerateNonceStr(),
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		SignType:  "MD5",
	}
	str, err = wxpay.NewTrade().Sign(defray, global.WechatApiKey)
	if err != nil {
		return wxpay.Defray{}, err
	}
	defray.PaySign = str
	return defray, nil
}

func generateTimestampStr() string {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func TestQrCode(t *testing.T) {
	// Create the barcode
	qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)

	// Scale the barcode to 200x200 pixels
	qrCode, _ = barcode.Scale(qrCode, 300, 300)
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	path := fmt.Sprintf("%s%s%v.png", global.FilesPath, "qrcode/test/", id)
	// create the output file

	file, flog := utils.CreateFile(path)

	defer file.Close()

	println(flog)
	// encode the barcode as png
	err := png.Encode(file, qrCode)
	if err != nil {
		t.Error(err)
	}
	println(qrCode.Content())
}

func TestFilesDir(t *testing.T) {
	path := fmt.Sprintf("%s%s", global.FilesPath, "qrcode/alipay")
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = strings.Replace(path, "/", "\\", -1)
	}
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		t.Error(err)
	}
}

func TestDirExists(t *testing.T) {
	path := fmt.Sprintf("%s%s", global.FilesPath, "qrcode/alipay")
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = strings.Replace(path, "/", "\\", -1)
	}
	flog, err := utils.PathExists(path)
	if err != nil {
		t.Error(err)
	}
	t.Log(flog)
}

func TestSubString(t *testing.T) {
	path := fmt.Sprintf("%s%s", global.FilesPath, "qrcode/alipay")
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = strings.Replace(path, "/", "\\", -1)
	}
	println(path)
	i := strings.LastIndex(path, "/")
	println(i)
	newPath := string([]rune(path)[0:i])
	t.Log(newPath)
}

func TestCreateFile(t *testing.T) {
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()

	path := fmt.Sprintf("%s%s%v.png", global.FilesPath, "qrcode/alipay/", id)
	println(path)
	_, flog := utils.CreateFile(path)
	println(flog)
}

func TestCreateQrCode(t *testing.T) {
	qrCode, _ := qr.Encode("测试", qr.M, qr.Auto)

	// Scale the barcode to 300*300 pixels
	qrCode, _ = barcode.Scale(qrCode, 300, 300)

	//iw, _ := utils.NewIdWorker(1)
	//id, _ := iw.NextId()
	//
	//path := fmt.Sprintf("%s%s%v.png", global.FilesPath, "qrcode/weichat/", id)
	//
	//// create the output file
	//file, err := utils.CreateFile(path)
	//if err != nil {
	//	return "", err
	//}
	//defer file.Close()
	//var file *os.File
	//// encode the barcode as png
	//png.Encode(file, qrCode)
	//
	//m, _ := png.Decode(file)
	//file.Close()

	// Do any adjustments to the image you need to
	// here such as resize and any filters
	// you might apply to the image

	var b bytes.Buffer
	png.Encode(&b, qrCode)

	img := imgbase64.FromBuffer(b)
	fmt.Println(img)
}

func TestFloat(t *testing.T) {
	f := convert.MustFloat64("1008") / 100
	println(convert.ToString(f))
}
