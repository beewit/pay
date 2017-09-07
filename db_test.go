package main

import (
	"encoding/json"
	"testing"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/pay/global"
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
	tx, _ := global.DB.Begin()
	for i := 10; i <= 20; i++ {
		x, err := tx.Insert("insert system_logs (id)values(?)", i)
		if err != nil {
			global.Log.Error("i:%v,err:%v", i, err.Error())
		}
		global.Log.Warning("i:%v,err:%v", "添加返回结果：", x)
	}
	err := tx.Commit()
	if err != nil {
		t.Error("TXCommit：" + err.Error())
	}
}
