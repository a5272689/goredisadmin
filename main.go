package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"github.com/goincremental/negroni-sessions"
	"goredisadmin/routers"
	"github.com/goincremental/negroni-sessions/redisstore"
	"net/http"
	"goredisadmin/models"
	"goredisadmin/modules"
	"goredisadmin/utils"
)




func main() {
	//defer models.Redis.Client.Close()
	r:= routers.Urls()
	n := negroni.Classic()
	rac,rc,cc:=utils.Rac,utils.Rc,utils.Cc
	store,err:=redisstore.New(20,"tcp",fmt.Sprintf("%v:%v",rc.Host,rc.Port),rc.Passwd,[]byte("secret123"))
	if err!=nil{
		fmt.Println(err)
	}
	sessionsH:=sessions.Sessions("my_session", store)
	userauth:=new(modules.AuthInfo)
  	n.Use(sessionsH)
	if cc.CasUrl!=""{
		casauth:=new(modules.CasAuthInfo)
		casauth.CasUrl=cc.CasUrl
		casauth.RedirectPath=cc.RedirectPath
		casauth.UserInfoApi=cc.UserInfoApi
		casauth.OpenAuth=cc.OpenAuth
		n.Use(casauth)
	}
	n.Use(userauth)
	n.Use(negroni.NewStatic(http.Dir(".")))
	n.UseHandler(r)
	listenaddr:=fmt.Sprintf("%v:%v",rac.Listen,rac.Port)
	go models.CheckRedis()
	n.Run(listenaddr)
}

