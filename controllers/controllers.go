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
	"goredisadmin/utils"
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
	hostname:=r.PostForm.Get("hostname")
	port,_:=strconv.Atoi(r.PostForm.Get("port"))
	utils.Logger.Printf("[info] SentinelsDataChangeAPI 收到参数：hostname:%v,port:%v",hostname,port)
	sentinel:=&models.Sentinel{Hostname:hostname,Port:port}
	saveresult,err:=sentinel.Create()
	result.Result=saveresult
	result.Info=fmt.Sprintf("报错：%v",err)
	jsonresult,_:=json.Marshal(result)
	strjsonresult:=string(jsonresult)
	utils.Logger.Printf("[info] SentinelsDataChangeAPI 结果：%v",strjsonresult)
	fmt.Fprint(w,strjsonresult)
}

func SentinelsDataDelAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	result:=new(JsonResult)
	result.Result=true
	data, _ := ioutil.ReadAll(r.Body)
	utils.Logger.Println("[info] SentinelsDataDelAPI 收到json串：",string(data))
	defer r.Body.Close()
	var del_sentinels []models.Sentinel
	json.Unmarshal(data,&del_sentinels)
	for _,tmp_sentinel_c:=range del_sentinels{
		tmp_del_result,_:=tmp_sentinel_c.Del()
		utils.Logger.Println("[info] SentinelsDataDelAPI 删除：",tmp_sentinel_c," 结果：",tmp_del_result)
	}
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
	AuthStatus bool `json:"auth_status"`
	PingStatus bool `json:"ping_status"`
	Version string `json:"version"`
	Role string `json:"role"`
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
			redisinfoslist=append(redisinfoslist,models.RedisInfo{Hostname:redislist[0],Port:redisport})
		}
	}
	alldata:=new(bootstrapTableRedissData)
	alldata.Rows=[]redissData{}

	for _,redisinfo:=range models.GetRediss(redisinfoslist...){
		alldata.Rows=append(alldata.Rows,redissData{Id:redisinfo["id"].(int),Version:redisinfo["version"].(string),
			Hostname:redisinfo["hostname"].(string),Port:redisinfo["port"].(int),AuthStatus:redisinfo["auth_status"].(bool),
			ConnectionStatus:redisinfo["connection_status"].(bool),PingStatus:redisinfo["ping_status"].(bool),
			Role:redisinfo["role"].(string),
		})
	}
	alldata.Total=len(alldata.Rows)
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
}

func RedissDataChangeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	r.ParseForm()
	result:=new(JsonResult)
	hostname:=r.PostForm.Get("hostname")
	port,_:=strconv.Atoi(r.PostForm.Get("port"))
	password:=r.PostForm.Get("password")
	utils.Logger.Printf("[info] RedissDataChangeAPI 收到参数：hostname:%v,port:%v,password:%v",hostname,port,password)
	redis:=&models.RedisInfo{Hostname:hostname,Port:port,Password:password}
	saveresult,err:=redis.Save()
	result.Result=saveresult
	result.Info=fmt.Sprintf("报错：%v",err)
	jsonresult,_:=json.Marshal(result)
	strjsonresult:=string(jsonresult)
	utils.Logger.Printf("[info] RedissDataChangeAPI 结果：%v",strjsonresult)
	fmt.Fprint(w,strjsonresult)
}

func RedissDataDelAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type","application/json")
	result:=new(JsonResult)
	result.Result=true
	data, _ := ioutil.ReadAll(r.Body)
	utils.Logger.Println("[info] reddissDataDelAPI 收到json串：",string(data))
	defer r.Body.Close()
	var del_rediss []models.RedisInfo
	json.Unmarshal(data,&del_rediss)
	for _,tmp_redis_c:=range del_rediss{
		tmp_del_result,_:=tmp_redis_c.Del()
		utils.Logger.Println("[info] reddissDataDelAPI 删除：",tmp_redis_c," 结果：",tmp_del_result)
	}
	jsonresult,_:=json.Marshal(result)
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