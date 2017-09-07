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
)
