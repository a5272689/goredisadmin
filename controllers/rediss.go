package controllers

import (
	"net/http"
	"github.com/flosch/pongo2"
	"goredisadmin/models"
	"github.com/goincremental/negroni-sessions"
	"io/ioutil"
	"goredisadmin/utils"
	"fmt"
	"github.com/bitly/go-simplejson"
	"strings"
	"strconv"
	"encoding/json"
)

func Rediss(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	redissstr:=strings.Join(r.Form["rediss"],"^")
	tpl,err:=pongo2.FromFile("views/contents/rediss.html")
	tpl = pongo2.Must(tpl,err)
	context:=initconText(r)
	context.Update(pongo2.Context{"redissstr":redissstr,"hiddenmastername":r.Form.Get("mastername")})
	tpl.ExecuteWriter(context, w)
}

type bootstrapTableRedissData struct {
	Rows []models.RedisInfo `json:"rows"`
	Total int `json:"total"`
}



func RedissDataAPI(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	jsonob,_:=simplejson.NewJson(data)
	redissstr,_:=jsonob.Get("rediss").String()
	mastername,_:=jsonob.Get("mastername").String()
	redisslist:=strings.Split(redissstr,"^")
	redisinfoslist:=[]models.RedisInfo{}
	for _,redisstr:=range redisslist{
		redislist:=strings.Split(redisstr,":")
		if len(redislist)==2{
			redisport,err:=strconv.Atoi(redislist[1])
			if err!=nil{
				continue
			}
			if mastername!=""{
				redisinfoslist=append(redisinfoslist,models.RedisInfo{Hostname:redislist[0],Port:redisport,Mastername:mastername})
			}else {
				redisinfoslist=append(redisinfoslist,models.RedisInfo{Hostname:redislist[0],Port:redisport})
			}

		}
	}
	alldata:=new(bootstrapTableRedissData)
	alldata.Rows=models.GetRediss(redisinfoslist...)
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
	mastername:=r.PostForm.Get("mastername")
	group:=r.PostForm.Get("group")
	savetype:=r.PostForm.Get("savetype")
	remark:=r.PostForm.Get("remark")
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] RedissDataChangeAPI 收到参数：hostname:%v,port:%v,password:%v,mastername:%v,group:%v,remark:%v,操作用户:%v",hostname,port,password,mastername,group,remark,session.Get("user"))
	redis:=&models.RedisInfo{Hostname:hostname,Port:port,Mastername:mastername,Password:password,Group:group,Remark:remark}
	var err error
	if savetype=="changepassword"{
		result.Result,err=redis.ChangePassword()
	}else if savetype=="change"{
		result.Result,err=redis.Change()
	} else {
		result.Result,err=redis.Save()
		go models.UpdateRedisInfo(hostname,port,models.GetHashName(hostname,port),mastername)
	}
	if err!=nil{
		utils.Logger.Println("[info] 保存报错：",err)
	}
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
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] reddissDataDelAPI 收到json串：%v,操作用户:%v",string(data),session.Get("user"))
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