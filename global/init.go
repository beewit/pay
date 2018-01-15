package global

import (
	"encoding/json"
	"fmt"

	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"github.com/beewit/beekit/utils/convert"
)

var (
	CFG  = conf.New("config.json")
	Log  = log.Logger
	DB   = mysql.DB
	DBTx *mysql.SqlConnTransaction
	RD   = redis.Cache
	IP   = CFG.Get("server.ip")
	Port = CFG.Get("server.port")
	Host = fmt.Sprintf("http://%v:%v", IP, Port)

	AlipayAppId         = fmt.Sprintf("%v", CFG.Get("alipay.appId"))
	AlipayPrivatePath   = fmt.Sprintf("%v", CFG.Get("alipay.privatePath"))
	AlipayPublicPath    = fmt.Sprintf("%v", CFG.Get("alipay.publicPath"))
	AlipayAliPublicPath = fmt.Sprintf("%v", CFG.Get("alipay.aliPublicPath"))
	AlipayReturnURL     = fmt.Sprintf("%v", CFG.Get("alipay.returnURL"))
	AlipayNotifyURL     = fmt.Sprintf("%v", CFG.Get("alipay.notifyURL"))

	WechatAppId     = fmt.Sprintf("%v", CFG.Get("wechat.appId"))
	WechatMchID     = fmt.Sprintf("%v", CFG.Get("wechat.mchID"))
	WechatApiKey    = fmt.Sprintf("%v", CFG.Get("wechat.apiKey"))
	WechatNotifyURL = fmt.Sprintf("%v", CFG.Get("wechat.notifyURL"))

	WechatAPPAppId     = fmt.Sprintf("%v", CFG.Get("wechatAPP.appId"))
	WechatAPPMchID     = fmt.Sprintf("%v", CFG.Get("wechatAPP.mchID"))
	WechatAPPApiKey    = fmt.Sprintf("%v", CFG.Get("wechatAPP.apiKey"))
	WechatAPPNotifyURL = fmt.Sprintf("%v", CFG.Get("wechatAPP.notifyURL"))

	WechatMiniAppConf = &wechatMiniAppConf{
		AppID:     convert.ToString(CFG.Get("wechat_mini_app.appId")),
		AppSecret: convert.ToString(CFG.Get("wechat_mini_app.appSecret")),
	}

	FilesPath   = fmt.Sprintf("%v", CFG.Get("files.path"))
	FilesDoMain = fmt.Sprintf("%v", CFG.Get("files.doMain"))
)

type wechatMiniAppConf struct {
	AppID     string
	AppSecret string
}
type Account struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Photo    string `json:"photo"`
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
}

func ToByteAccount(b []byte) *Account {
	var rp = new(Account)
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	return rp
}

func ToMapAccount(m map[string]interface{}) *Account {
	b := convert.ToMapByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}

func ToInterfaceAccount(m interface{}) *Account {
	b := convert.ToInterfaceByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}
