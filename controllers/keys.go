package controllers

import (
	"fmt"
	"encoding/json"
	"github.com/flosch/pongo2"
	"net/http"
	"github.com/goincremental/negroni-sessions"
	//"fmt"
	//"time"
	//"encoding/json"
	"goredisadmin/models"
	"goredisadmin/utils"
	//"github.com/bitly/go-simplejson"
	//"strconv"
	//"strings"
	"io/ioutil"
	//"goredisadmin/utils"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
)

func Keys(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	fmt.Println(session.Get("user"))
	redis_list:=[]string{}
	r.ParseForm()
	redis:=r.Form.Get("redis")
	if redis!=""{
		redis_list=append(redis_list,redis)
	}else {
		redis_list=models.GetRedisNames()
	}
	dbs_json, _ := json.Marshal(models.GetRedisDbs(redis_list))
	tpl,err:=pongo2.FromFile("views/contents/keys.html")
	tpl = pongo2.Must(tpl,err)
	context:=initconText(r)
	context=context.Update(pongo2.Context{"rediss":redis_list,"db_map":string(dbs_json)})
	tpl.ExecuteWriter(context, w)
}





type bootstrapTableKeysData struct {
	Rows []models.KeysData `json:"rows"`
	Total int `json:"total"`
}


func KeysDataAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	utils.Logger.Println(string(data))
	jsonob,_:=simplejson.NewJson(data)
	keysstr,_:=jsonob.Get("keys").String()
	redisstr,_:=jsonob.Get("redis").String()
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	roleinfo:=redisinfo.GetRoleInfo()
	if roleinfo.Role=="master"&&len(roleinfo.Slaves)>0{
		redisinfo=*roleinfo.Slaves[0]
	}
	alldata:=new(bootstrapTableKeysData)
	session := sessions.GetSession(r)
	userrole:=session.Get("role")
	utils.Logger.Println(userrole,keysstr,redis_db_index)
	if userrole=="ops"||keysstr!="*"{
		alldata.Rows=redisinfo.GetKeys(keysstr,redis_db_index)
	}
	alldata.Total=len(alldata.Rows)
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
}

