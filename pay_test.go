package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"image/png"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/imgbase64"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/pay/alipay"
	"github.com/beewit/pay/global"
	"github.com/beewit/pay/handler"
	"github.com/beewit/pay/wxpay"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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
	}
	str, err = wxpay.NewTrade().Sign(defray, global.WechatApiKey)
	if err != nil {
		return wxpay.Defray{}, err
	}
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
	path := fmt.Sprintf("%s%s%v.png", global.FilesPath, "qrcode/test/", utils.ID())
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
	path := fmt.Sprintf("%s%s%v.png", global.FilesPath, "qrcode/alipay/", utils.ID())
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

func TestAddDay(t *testing.T) {
	expirTimeStr, _ := time.Parse("2006-01-02 15:04:05", "2017-11-29 00:00:00")
	println(utils.FormatTime(expirTimeStr))
	println(utils.FormatTime(expirTimeStr.AddDate(0, 0, 90)))
	println(utils.FormatTime(expirTimeStr.AddDate(0, 0, 95)))
}

func TestNofily(t *testing.T) {
	order := handler.GetOrder(6433434672055296)
	flog := handler.UpdateOrderFuncStatus(order, 10, "101.27.240.120")
	println(flog)
}

func TestTxInsertMap(t *testing.T) {
	flog := false
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		m := map[string]interface{}{}
		m["id"] = 123456789
		m["account_id"] = 123456789
		m["func_id"] = 501
		m["expiration_time"] = utils.CurrentTime()
		m["ct_time"] = utils.CurrentTime()
		m["ut_time"] = utils.CurrentTime()
		_, err := tx.InsertMap("account_func", m)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		flog = true
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
			flog = false
		}
	})
	println(flog)
}

func TestFor(t *testing.T) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i == j {
				println("跳出循环 - 》 I：", i, "J：", j)
				break
			}
			println("I：", i, "J：", j)
		}
	}
}

func TestCreateOrder(t *testing.T) {

	println(fmt.Sprintf("%.0f", 51.9*0.1))
}

func createOrder() {
	m := make(map[string]interface{})
	tradeNo := utils.ID()
	//测试使用
	//totalPrice = 0.01

	m["id"] = tradeNo
	m["account_id"] = 1
	m["account_name"] = "工蜂小智"
	m["type"] = enum.PAY_TYPE_FUNC
	m["pay_type"] = "微信"
	m["pay_price"] = 0.01
	m["pay_status"] = enum.PAY_STATUS_NOT
	m["status"] = enum.NORMAL
	m["ct_time"] = utils.CurrentTime()
	m["ct_ip"] = "123"
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		_, err := tx.InsertMap("order_payment", m)
		if err != nil {
			panic(err)
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)

		}
	})
}

func TestRedPacketPay(t *testing.T) {
	e := echo.New()
	f := url.Values{}
	f.Set("id", "1")
	f.Set("openId", "oNL8m0T-FAk0NadV4rktzT5G-na4")
	f.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Jis3756dX_49PhACfn9xl_Fn0MzWJI0Hyewb9my3SxY")
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	m, _ := global.DB.Query("SELECT id,nickname,photo,mobile,status FROM account WHERE id=? LIMIT 1", 122068319091036160)
	c.Set("account", global.ToMapAccount(m[0]))

	// 断言
	if assert.NoError(t, handler.RedPacketPay(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t.Error(rec.Body.String())
	}
}

func TestNotityRedPacketPay(t *testing.T) {
	order := handler.GetOrder(6426687589467136)
	handler.UpdateOrderRedPacketStatus(order, 1, "127.0.0.1")
}

func TestQrcode(t *testing.T) {
	body, _ := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    fmt.Sprintf("http://m.9ee3.com/account/create/temporary/qrcode?objId=%v&objType=%s", 1, enum.QRCODE_RED_PACKET),
	})
	t.Error(string(body))
}
