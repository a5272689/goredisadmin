package utils

import (
	"github.com/go-ini/ini"
	"path/filepath"
	"os"
	"fmt"
)

type RedisAdminConf struct {
	Listen string `ini:"listen"`
	Port  int `ini:"port"`
}
type RedisConf struct {
	Host string `ini:"host"`
	Port  int `ini:"port"`
	Passwd string `ini:"passwd"`
}
type CasConf struct {
	CasUrl string `ini:"casurl"`
	RedirectPath  string `ini:"redirectpath"`
	UserInfoApi string `ini:"userinfoapi"`
	OpenAuth bool `ini:"openauth"`
}
func ConfLoad() (*RedisAdminConf,*RedisConf,*CasConf) {
	rac := new(RedisAdminConf)
	rac.Port=3000
	rac.Listen="0.0.0.0"
	rc := new(RedisConf)
	rc.Port=6379
	rc.Host="127.0.0.1"
	cc:=new(CasConf)
	cc.RedirectPath="/"
	cc.OpenAuth=false
	logger:=Logger
	selfpath,_:=filepath.Abs(os.Args[0])
	basedir,_:=filepath.Split(selfpath)
	defaultConf:=filepath.Join(basedir,"conf","goredisadmin.ini")
	//defaultConf="/root/Desktop/gopath/src/goredisadmin/conf/goredisadmin.ini"
	cfg, err := ini.Load(defaultConf)
	if err!=nil{
		logger.Println(fmt.Sprintf("[info] 加载配置文件：%v失败！！！",defaultConf))
		return rac,rc,cc
	}
	cfg.MapTo(rac)
	cfg.Section("Redis").MapTo(rc)
	cfg.Section("Cas").MapTo(cc)
	return rac,rc,cc
}

var Rac,Rc,Cc=ConfLoad()