package modules

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/urfave/negroni"
)

func NewAuth() negroni.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		byteurl:=[]byte(r.URL.Path)
		bytestatic:=[]byte("/static")
		bytelogin:=[]byte("/login")
		session := sessions.GetSession(r)
		if string(byteurl[0:len(bytelogin)])!="/login" && string(byteurl[0:len(bytestatic)])!="/static" &&session.Get("user")==nil{
			http.Redirect(w,r,"/login",http.StatusFound)
		}else {
			next(w, r)
		}
	}
}


