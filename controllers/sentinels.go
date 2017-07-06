package controllers

import (
	"github.com/flosch/pongo2"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/goincremental/negroni-sessions"
	"goredisadmin/models"
	"strconv"
	"goredisadmin/utils"
	"io/ioutil"
)

func Sentinels(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	tpl,err:=pongo2.FromFile("views/contents/sentinels.html")
	tpl = pongo2.Must(tpl,err)
	fmt.Println(session.Get("user"))
	tpl.ExecuteWriter(initconText(r), w)
}

type bootstrapTableSentinelsData struct {
	Rows []models.Sentinel `json:"rows"`
	Total int `json:"total"`
}



func SentinelsDataAPI(w http.ResponseWriter, r *http.Request) {
	alldata:=new(bootstrapTableSentinelsData)
	alldata.Rows,_=models.GetSentinels()
	alldata.Total=len(alldata.Rows)
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
}





func SentinelsDataChangeAPI(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	w.Header().Add("Content-Type","application/json")
	r.ParseForm()
	result:=new(JsonResult)
	hostname:=r.PostForm.Get("hostname")
	port,_:=strconv.Atoi(r.PostForm.Get("port"))
	utils.Logger.Printf("[info] SentinelsDataChangeAPI 收到参数：hostname:%v,port:%v,操作用户:%v",hostname,port,session.Get("user"))
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
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] SentinelsDataDelAPI 收到json串：%v,操作用户:%v",string(data),session.Get("user"))
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


