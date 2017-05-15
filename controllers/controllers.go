package controllers

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"fmt"
	"github.com/flosch/pongo2"
	"time"
	"encoding/json"
	"crypto/sha256"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Set("hello", "world2")
	http.SetCookie(w,&http.Cookie{Name:"csrftoken",Value:string(time.Now().String()),MaxAge:60})
	fmt.Println(ConfLoad())
	fmt.Fprintln(w, "a")
	//http.Redirect(w,r,"/",http.StatusFound)

}


func Main2Handler(w http.ResponseWriter, r *http.Request) {

	session := sessions.GetSession(r)
	tpl,err:=pongo2.FromFile("views/abcd.html")
	tplExample := pongo2.Must(tpl,err)
	abc,_:=tplExample.ExecuteBytes(pongo2.Context{"user": "hjd","qq":"123","user2":session.Get("hello")})
	for _,c:=range r.Cookies(){
		fmt.Println(c.Name,c.Value)
	}

	//session.Set("hello", "world")
	//http.Redirect(w,r,"/a",http.StatusFound)
	fmt.Fprintln(w, string(abc))

	//fmt.Println(Redis.Cmd("set","goredisadmin:user:abc","av"))

	//http.ResponseWriter()
}

func Login(w http.ResponseWriter, r *http.Request)  {
	session := sessions.GetSession(r)
	user:=session.Get("user")
	if user==nil{
		tpl,err:=pongo2.FromFile("views/login.html")
		tpl = pongo2.Must(tpl,err)
		tpl.ExecuteWriter(pongo2.Context{"title":"Redis Admin Login"}, w)
	}else {
		http.Redirect(w,r,"/",http.StatusFound)
	}
}

type LoginResult struct {
	Result bool `json:"result"`
	Info string `json:"info"`
}

func LoginAuth(w http.ResponseWriter, r *http.Request)  {
	w.Header().Add("Content-Type","application/json")
	session := sessions.GetSession(r)
	result:=new(LoginResult)
	r.ParseForm()
	username:=r.PostForm.Get("username")
	passwd:=r.PostForm.Get("passwd")
	dbpass,_:=Redis.Cmd("get","goredisadmin:user:"+username).Str()
	h:=sha256.New()
	h.Write([]byte(passwd))
	h.Write([]byte(string(len(passwd))))
	h.Write([]byte("goredisadmin"))
	tmp_passwd:=fmt.Sprintf("%x", h.Sum(nil))
	if dbpass==tmp_passwd && len(passwd)>0{
		session.Set("user",username)
		result.Result=true
	}else {
		result.Info=fmt.Sprintf("用户%v认证失败！！！",username)
	}
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}