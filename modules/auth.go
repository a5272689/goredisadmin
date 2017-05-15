package modules

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"github.com/urfave/negroni"
	"strings"
)

func NewAuth() negroni.HandlerFunc  {
	return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		urlpsth:=r.URL.Path
		session := sessions.GetSession(r)
		if strings.HasPrefix(urlpsth,"/login") || strings.HasPrefix(urlpsth,"/static") || session.Get("user")!=nil{
			next(w, r)
		}else {
			http.Redirect(w,r,"/login",http.StatusFound)
		}
	}
}


