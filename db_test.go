package main

import (
	"encoding/json"
	"testing"
	"encoding/base64"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/pay/global"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils/convert"
	"time"
	"github.com/beewit/pay/handler"
	"net/http"
	"encoding/xml"
	"github.com/beewit/pay/wxpay"
)

func TestMemberType(t *testing.T) {

	b, err := uhttp.PostForm("http://127.0.0.1:8083/member/type", nil)
	if err != nil {
		global.Log.Error(err.Error())
		t.Error(err)
	}
	println(string(b[:]))
	var rp utils.ResultParam
	json.Unmarshal(b[:], &rp)
	println("牛逼了")
}

func TestTx(t *testing.T) {
	errCount := 0
	tx, _ := global.DB.Begin()
	for i := 10; i <= 20; i++ {
		x, err := tx.Insert("insert system_logs (id)values(?)", i)
		if err != nil {
			errCount++
			global.Log.Error("跳出循环：i:%v,err:%v", i, err.Error())
		}
		global.Log.Warning("i:%v,err:%v", "添加返回结果：", x)
	}
	if errCount > 0 {
		err := tx.Rollback()
		if err != nil {
			t.Error("TXRollback：" + err.Error())
		}
	} else {
		err := tx.Commit()
		if err != nil {
			t.Error("TXCommit：" + err.Error())
		}
	}

}

func TestTxCommon(t *testing.T) {
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		for i := 22; i <= 30; i++ {
			x, err := tx.Insert("insert system_logs (id)values(?)", i)
			if err != nil {
				panic(err)
				global.Log.Error("跳出循环：i:%v,err:%v", i, err.Error())
			}
			global.Log.Warning("i:%v,err:%v", "添加返回结果：", x)
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
		}
	})
}

func TestTxCommonMap(t *testing.T) {
	flog := true
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		for i := 43; i <= 44; i++ {
			m := make(map[string]interface{})
			m["id"] = i
			m["content"] = convert.ToString(i) + "张三很累了啊"
			x, err := global.DB.InsertMap("system_logs", m)
			if err != nil {
				panic(err)
			} else {
				t.Log(convert.ToString(i), "结果： ", x)
			}
		}
	}, func(err error) {
		if err != nil {
			global.Log.Error("保存失败，%v", err)
			flog = false
		}
	})
	if !flog {
		global.Log.Error("保存失败")
	} else {
		global.Log.Error("保存成功222222222")
	}
}

func TestInsertMap(t *testing.T) {
	m := make(map[string]interface{})
	m["id"] = 2
	m["content"] = "张三很累了啊"
	x, err := global.DB.InsertMap("system_logs", m)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("结果： ", x)
	}
}

func TestGetIp(t *testing.T) {
	println(utils.GetIp())
}

func TestBase64(t *testing.T) {
	str := `OTJLBQZ7YjSgmRL5gP/x5gMHOSzawJScxVM5FVpX3RRgMaVrJkXmMsFX6lnxfPTQbVLqY9WeUJ4nCp/JWdPzjVdXwSjVacGfvdaO92mAsiiGJ0yBaMJ8tbdNVlxnLHA0PNpwv2K+RINchWgFxLQ+e6Qpe60r9VdBhMNW6tKLfrfATLQU0eLLjJXzY0R01Idtqwri06BYjPpClrcxwykPkEJDM1OZxcIxm7+pcOQ8VdgHoj9krVjrDcV7J6hPccXPRh1MgPaWm9e5nqkf60PS3BHqlbt5z189PcVxood7km+YX4+rKuHBzz8VudwyHLn1Zx3RhF3zTSGi59hJpu1UBQ==`
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		t.Error(err)
	}
	t.Log(b)
}

func TestConvert(t *testing.T) {
	var id int64
	id = 123456789
	println(convert.ToString(id))

	tim := time.Now().AddDate(0, 0, 365)
	println(utils.FormatTime(tim))
}

func TestUpdateOrder(t *testing.T) {
	flog := handler.UpdateOrderMTCStatus(125985189636608000, 0.1)
	println(flog)
}

func TestTime(t *testing.T) {
	iw, _ := utils.NewIdWorker(1)
	tradeNo, err := iw.NextId()
	if err != nil {
		t.Error(err)
	}
	println(convert.ToString(tradeNo))
	println(len(convert.ToString(tradeNo)))
}

func TestID(t *testing.T) {
	iw, _ := utils.NewIdWorker(1)
	tradeNo, _ := iw.NextId()
	tradeNo2, _ := iw.NextId()
	tradeNo3, _ := iw.NextId()
	println(tradeNo)
	println(tradeNo2)
	println(tradeNo3)
}

func TestWechatNotify(t *testing.T) {
	wn := &wxpay.Notice{
		AppID: "123456789",
	}
	body, err := xml.Marshal(wn)
	if err != nil {
		t.Error(err)
	}
	header := http.Header{}
	header.Add("Accept", "application/xml")
	header.Add("Content-Type", "application/xml;charset=utf-8")
	body, err = uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    "http://127.0.0.1:8083/wechat/notify",
		Body:   body,
		Header: header,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body[:]))
}
