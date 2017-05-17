package modules

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"net/url"
	"io/ioutil"
	"encoding/xml"
)

type CasAuthInfo struct {
	CasUrl string
	RedirectPath string
	UserName string
}

type UserXml struct {
	User string `xml:"authenticationSuccess>user"`
}

func (c *CasAuthInfo) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)()  {
	urlpsth:=r.URL.Path
	session := sessions.GetSession(r)
	r.ParseForm()
	ticket:=r.Form.Get("ticket")
	if ticket!=""{
		a:=url.Values{"service":{"http://"+r.Host+c.RedirectPath},"ticket":{ticket},"format":{"JSON"}}
		resp,_:=http.Get(c.CasUrl+"/p3/serviceValidate?"+a.Encode())
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		userxml:=&UserXml{}
		xml.Unmarshal(body,userxml)
		session.Set("casuser",userxml.User)
		http.Redirect(rw,r,"http://"+r.Host+"/login",http.StatusFound)
	}else {
		if urlpsth=="/login" && session.Get("casuser")==nil{
			a:=url.Values{"service":{"http://"+r.Host+c.RedirectPath}}
			http.Redirect(rw,r,c.CasUrl+"/login"+"?"+a.Encode(),http.StatusFound)
		}else if  urlpsth=="/logout"{
			session.Clear()
			a:=url.Values{"service":{"http://"+r.Host+c.RedirectPath}}
			http.Redirect(rw,r,c.CasUrl+"/logout"+"?"+a.Encode(),http.StatusFound)
		}else {
			next(rw,r)
		}
	}

}
