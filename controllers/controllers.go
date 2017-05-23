package controllers

import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"fmt"
	//"time"
	"encoding/json"
	"github.com/flosch/pongo2"
	"goredisadmin/models"
	//"goredisadmin/utils"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
	"io/ioutil"
)

func initconText(r *http.Request) pongo2.Context {
	session := sessions.GetSession(r)
	return pongo2.Context{"username":session.Get("user"),"urlpath":r.URL.Path}
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	//http.SetCookie(w,&http.Cookie{Name:"csrftoken",Value:string(time.Now().String()),MaxAge:60})
	//fmt.Println(r.URL.Path)
	//fmt.Println(utils.ConfLoad())
	//fmt.Fprintln(w, session.Get("user"))
	//userdb:=&models.User{UserName:"jkljdaklsjfkl"}
	//dbpass,err:=userdb.GetPassWord()
	//fmt.Println(dbpass,err)
	//http.Redirect(w,r,"/",http.StatusFound)
	tpl,err:=pongo2.FromFile("views/contents/index.html")
	tpl = pongo2.Must(tpl,err)
	fmt.Println(session.Get("user"))
	tpl.ExecuteWriter(initconText(r), w)
}

func Sentinels(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	tpl,err:=pongo2.FromFile("views/contents/sentinels.html")
	tpl = pongo2.Must(tpl,err)
	fmt.Println(session.Get("user"))
	tpl.ExecuteWriter(initconText(r), w)
}

type bootstrapTableSentinelsData struct {
	Rows []sentinelsData `json:"rows"`
	Total int `json:"total"`
}

type sentinelsData struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Masters []string `json:"masters"`
	ConnectionStatus bool `json:"connection_status"`
	MasterRediss map[string][]map[string]string `json:"master_rediss"`
	Version string `json:"version"`
	
}

func SentinelsDataAPI(w http.ResponseWriter, r *http.Request) {
	alldata:=new(bootstrapTableSentinelsData)
	alldata.Rows=[]sentinelsData{}
	for _,sentinel:=range models.GetSentinels(){
		alldata.Rows=append(alldata.Rows,sentinelsData{Id:sentinel["id"].(int),Version:sentinel["version"].(string),
			Hostname:sentinel["hostname"].(string),Port:sentinel["port"].(int),Masters:sentinel["masters"].([]string),
			ConnectionStatus:sentinel["connection_status"].(bool),MasterRediss:sentinel["master_rediss"].(map[string][]map[string]string),
		})
	}
	alldata.Total=len(alldata.Rows)
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
}

func SentinelsDataChangeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	r.ParseForm()
	result:=new(JsonResult)
	sentinelid:=r.PostForm.Get("sentinelid")
	hostname:=r.PostForm.Get("hostname")
	port,_:=strconv.Atoi(r.PostForm.Get("port"))
	sentinel:=&models.Sentinel{HostName:hostname,Port:port,Sentinelid:sentinelid}
	saveresult,err:=sentinel.Save()
	result.Result=saveresult
	result.Info=fmt.Sprintf("报错：%v",err)
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}

func Rediss(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	r.ParseForm()
	redissstr:=strings.Join(r.Form["rediss"],"^")
	tpl,err:=pongo2.FromFile("views/contents/rediss.html")
	tpl = pongo2.Must(tpl,err)
	fmt.Println(session.Get("user"))
	context:=initconText(r)
	context.Update(pongo2.Context{"redissstr":redissstr})
	tpl.ExecuteWriter(context, w)
}

type bootstrapTableRedissData struct {
	Rows []redissData `json:"rows"`
	Total int `json:"total"`
}

type redissData struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	ConnectionStatus bool `json:"connection_status"`
	Version string `json:"version"`
}

func RedissDataAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	jsonob,_:=simplejson.NewJson(data)
	redissstr,_:=jsonob.Get("rediss").String()
	redisslist:=strings.Split(redissstr,"^")
	redisinfoslist:=[]models.RedisInfo{}
	for _,redisstr:=range redisslist{
		redislist:=strings.Split(redisstr,":")
		if len(redislist)==2{
			redisport,err:=strconv.Atoi(redislist[1])
			if err!=nil{
				continue
			}
			redisinfoslist=append(redisinfoslist,models.RedisInfo{HostName:redislist[0],Port:redisport})
		}
	}
	alldata:=new(bootstrapTableRedissData)
	alldata.Rows=[]redissData{}
	models.GetRediss(redisinfoslist...)
	//for _,_:=range models.GetRediss(){
	//	alldata.Rows=append(alldata.Rows,redissData{Id:sentinel["id"].(int),Version:sentinel["version"].(string),
	//		Hostname:sentinel["hostname"].(string),Port:sentinel["port"].(int),Masters:sentinel["masters"].([]string),
	//		ConnectionStatus:sentinel["connection_status"].(bool),MasterRediss:sentinel["master_rediss"].(map[string][]map[string]string),
	//	})
	//}
	alldata.Total=len(alldata.Rows)
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
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
	if user==nil{
		tpl,err:=pongo2.FromFile("views/login.html")
		tpl = pongo2.Must(tpl,err)
		tpl.ExecuteWriter(pongo2.Context{"title":"Redis Admin Login"}, w)
	} else {
		http.Redirect(w,r,"/",http.StatusFound)
	}
}

type JsonResult struct {
	Result bool `json:"result"`
	Info string `json:"info"`
}

func LoginAuth(w http.ResponseWriter, r *http.Request)  {
	w.Header().Add("Content-Type","application/json")
	session := sessions.GetSession(r)
	result:=new(JsonResult)
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