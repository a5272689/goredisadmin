package modules

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"strings"
)

type AuthInfo struct {
	UserName string
}

func (a *AuthInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)()  {
	urlpath:=r.URL.Path
	session := sessions.GetSession(r)
	if strings.HasPrefix(urlpath,"/login") || strings.HasPrefix(urlpath,"/static") || session.Get("user")!=nil{
		role:=session.Get("role")
		if role==nil{
			if urlpath=="/sentinelschange"|| urlpath=="/sentinelsdel"||urlpath=="/redisschange"||
				urlpath=="/redissdel"||urlpath=="/keysdel" ||urlpath=="/keysexpire"||urlpath=="/keyvaldel"||
				urlpath=="/keyspersist"||urlpath=="/keysave" ||urlpath=="/keyrename"{
				http.NotFound(rw,r)
			}else {
				next(rw, r)
			}
		}else {
			next(rw, r)
		}
	}else {
		http.Redirect(rw,r,"/login",http.StatusFound)
	}
}



