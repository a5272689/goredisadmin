package controllers

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"fmt"
	"time"
	"encoding/json"
	"github.com/flosch/pongo2"
	"goredisadmin/models"
	"goredisadmin/utils"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	http.SetCookie(w,&http.Cookie{Name:"csrftoken",Value:string(time.Now().String()),MaxAge:60})
	fmt.Println(r.URL.Path)
	fmt.Println(utils.ConfLoad())
	fmt.Fprintln(w, session.Get("user"))
	userdb:=&models.User{UserName:"jkljdaklsjfkl"}
	dbpass,err:=userdb.GetPassWord()
	fmt.Println(dbpass,err)
	//http.Redirect(w,r,"/",http.StatusFound)

}


func Logout(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	session.Clear()
	http.Redirect(w,r,"/",http.StatusFound)
}

func Login(w http.ResponseWriter, r *http.Request)  {
	session := sessions.GetSession(r)
	casuser:=session.Get("casuser")
	r.ParseForm()
	if casuser!=nil{
		session.Set("user",casuser)
	}
	user:=session.Get("user")
	fmt.Println(casuser,user)
	if user==nil{
		tpl,err:=pongo2.FromFile("views/login.html")
		tpl = pongo2.Must(tpl,err)
		tpl.ExecuteWriter(pongo2.Context{"title":"Redis Admin Login"}, w)
	} else {
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
	userinfo:=&models.User{UserName:username}
	dbpass,_:=userinfo.GetPassWord()
	tmp_passwd:=userinfo.HashPasswd(passwd)
	if dbpass==tmp_passwd && len(passwd)>0{
		session.Set("user",username)
		result.Result=true
	}else {
		result.Info=fmt.Sprintf("用户%v认证失败！！！",username)
	}
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}