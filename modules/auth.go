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
	urlpsth:=r.URL.Path
	session := sessions.GetSession(r)
	if strings.HasPrefix(urlpsth,"/login") || strings.HasPrefix(urlpsth,"/static") || session.Get("user")!=nil{
		next(rw, r)
	}else {
		http.Redirect(rw,r,"/login",http.StatusFound)
	}
}



