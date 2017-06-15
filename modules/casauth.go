package modules

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"net/url"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
)

type CasAuthInfo struct {
	CasUrl string
	RedirectPath string
	UserInfoApi string
	UserName string
	OpenAuth bool
}

type UserXml struct {
	User string `xml:"authenticationSuccess>user"`
}

type UserInfoJson struct {
	Data UserInfoNameJson `json:"data"`
}

type UserInfoNameJson struct {
	Name string `json:"name"`
	Depart UserInfoPNameJson `json:"depart"`
}

type UserInfoPNameJson struct {
	Name string `json:"name"`
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
			if c.UserInfoApi!=""&&session.Get("casuser")!=nil{
				res, _ := http.Get(c.UserInfoApi+"/"+session.Get("casuser").(string));
				result, _ := ioutil.ReadAll(res.Body)
				res.Body.Close()
				userinfo:=&UserInfoJson{}
				json.Unmarshal(result,userinfo)
				session.Set("username",userinfo.Data.Name)
				if userinfo.Data.Depart.Name=="基础运维"||c.OpenAuth{
					session.Set("role","ops")
				}
			}
			next(rw,r)
		}
	}

}
