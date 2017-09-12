package global

import (
	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"fmt"
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

	FilesPath   = fmt.Sprintf("%v", CFG.Get("files.path"))
	FilesDoMain = fmt.Sprintf("%v", CFG.Get("files.doMain"))
)
