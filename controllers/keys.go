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


func KeysDataDelAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeysDataDelAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	tmpkeyslist,_:=jsonob.Get("keys").Array()
	keyslist:=[]string{}
	for _,keyname:=range tmpkeyslist{
		keyslist=append(keyslist,keyname.(string))
	}
	redisstr,_:=jsonob.Get("redis").String()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	del_keys:=redisinfo.DelKeys(keyslist,redis_db_index)
	result:=new(JsonResult)
	result.Result=true
	result.Info=strings.Join(del_keys,",")
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}




func KeysDataExpireAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeysDataExpireAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	tmpkeyslist,_:=jsonob.Get("keys").Array()
	keyslist:=[]string{}
	for _,keyname:=range tmpkeyslist{
		keyslist=append(keyslist,keyname.(string))
	}
	redisstr,_:=jsonob.Get("redis").String()
	seconds,_:=jsonob.Get("seconds").Int()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	expire_keys:=redisinfo.ExpireKeys(keyslist,seconds,redis_db_index)
	result:=new(JsonResult)
	result.Result=true
	result.Info=strings.Join(expire_keys,",")
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}

func KeysDataPersistAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeysDataPersistAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	tmpkeyslist,_:=jsonob.Get("keys").Array()
	keyslist:=[]string{}
	for _,keyname:=range tmpkeyslist{
		keyslist=append(keyslist,keyname.(string))
	}
	redisstr,_:=jsonob.Get("redis").String()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	persist_keys:=redisinfo.PersistKeys(keyslist,redis_db_index)
	result:=new(JsonResult)
	result.Result=true
	result.Info=strings.Join(persist_keys,",")
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}



func KeyRenameAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeyRenameAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	key,_:=jsonob.Get("key").String()
	newkey,_:=jsonob.Get("newkey").String()
	redisstr,_:=jsonob.Get("redis").String()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	renamekeys,err:=redisinfo.RenameKey(key,newkey,redis_db_index)
	result:=new(JsonResult)
	if err==nil{
		if renamekeys==1{
			result.Result=true
		}
	}
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}



func KeyValDelAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeyValDelAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	key,_:=jsonob.Get("key").String()
	key_type,_:=jsonob.Get("type").String()
	redisstr,_:=jsonob.Get("redis").String()
	field,_:=jsonob.Get("field").String()
	member,_:=jsonob.Get("val").String()
	indexstr,_:=jsonob.Get("index").String()
	index,_:=strconv.Atoi(indexstr)
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	result:=new(JsonResult)
	result.Result=true
	switch key_type {
	case "string":
		_,err:=redisinfo.DelStrValKey(key,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "hash":
		_,err:=redisinfo.DelHashValKey(key,field,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "list":
		err:=redisinfo.DelListValKey(key,index,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "set":
		_,err:=redisinfo.DelSetValKey(key,member,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "zset":
		_,err:=redisinfo.DelZsetValKey(key,member,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	default:
		result.Result=false
		result.Info="不支持的类型"
	}
	utils.Logger.Println(key,key_type,redisport,redis_db_index)
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}

func KeySaveAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	session := sessions.GetSession(r)
	utils.Logger.Printf("[info] KeySaveAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
	jsonob,_:=simplejson.NewJson(data)
	key_type,_:=jsonob.Get("type").String()
	key,_:=jsonob.Get("key").String()
	val,_:=jsonob.Get("val").String()
	field,_:=jsonob.Get("field").String()
	indexstr,_:=jsonob.Get("index").String()
	scorestr,_:=jsonob.Get("score").String()
	score,_:=strconv.Atoi(scorestr)
	redisstr,_:=jsonob.Get("redis").String()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	result:=new(JsonResult)
	result.Result=true
	switch key_type {
	case "string":
		_,err:=redisinfo.SetKey(key,val,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "hash":
		_,err:=redisinfo.HsetKey(key,field,val,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "list":
		var err error
		if indexstr==""{
			_,err=redisinfo.LpushKey(key,val,redis_db_index)
		}else {
			index,_:=strconv.Atoi(indexstr)
			_,err=redisinfo.LsetKey(key,index,val,redis_db_index)
		}

		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "set":
		_,err:=redisinfo.SaddKey(key,val,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	case "zset":
		_,err:=redisinfo.ZaddKey(key,score,val,redis_db_index)
		if err!=nil{
			result.Result=false
			result.Info=fmt.Sprintf("报错：%v",err)
		}
	default:
		result.Result=false
		result.Info="不支持的类型"
	}
	jsonresult,_:=json.Marshal(result)
	fmt.Fprint(w,string(jsonresult))
}



func KeyDataAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	utils.Logger.Println(string(data))
	jsonob,_:=simplejson.NewJson(data)
	keystr,_:=jsonob.Get("key").String()
	redisstr,_:=jsonob.Get("redis").String()
	redislist:=strings.Split(redisstr,":")
	redisport,_:=strconv.Atoi(redislist[1])
	redis_db,_:=jsonob.Get("redis_db").String()
	redis_db_index,_:=strconv.Atoi(redis_db)
	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
	key_type,err:=redisinfo.TypeKey(keystr,redis_db_index)
	utils.Logger.Println(key_type,err)
	key_ttl,err:=redisinfo.TtlKey(keystr,redis_db_index)
	utils.Logger.Println(key_ttl,err)
	alldata:=make(map[string]interface{})
	rows:=[]map[string]string{}
	alldata["type"]=key_type
	alldata["ttl"]=key_ttl
	switch key_type {
	case "string":
		val,err:=redisinfo.GetKey(keystr,redis_db_index)
		utils.Logger.Println(val)
		valmap:=make(map[string]string)
		if err==nil{
			valmap["val"]=val
			rows=append(rows,valmap)

		}
	case "hash":
		vals,err:=redisinfo.HmgetKey(keystr,redis_db_index)
		if err==nil{
			for field,val:=range vals{
				valmap:=make(map[string]string)
				valmap["val"]=val
				valmap["field"]=field
				rows=append(rows,valmap)
			}
		}
		utils.Logger.Println(vals)
	case "list":
		vals,err:=redisinfo.LrangeKey(keystr,redis_db_index)
		if err==nil{
			for index,val:=range vals {
				valmap := make(map[string]string)
				valmap["val"] = val
				valmap["index"]=strconv.Itoa(index)
				rows = append(rows, valmap)
			}
		}
		utils.Logger.Println(vals)
	case "set":
		vals,err:=redisinfo.SmembersKey(keystr,redis_db_index)
		if err==nil{
			for index,val:=range vals {
				valmap := make(map[string]string)
				valmap["val"] = val
				valmap["index"]=strconv.Itoa(index)
				rows = append(rows, valmap)
			}
		}
		utils.Logger.Println(vals)
	case "zset":
		vals,err:=redisinfo.ZrangeKey(keystr,redis_db_index)
		if err==nil{
			for index,val:=range vals {
				if index%2==0{
					valmap := make(map[string]string)
					valmap["val"] = val
					valmap["score"]=vals[index+1]
					rows = append(rows, valmap)
				}
			}
		}
	}
	alldata["rows"]=rows
	jsonresult,_:=json.Marshal(alldata)
	fmt.Fprint(w,string(jsonresult))
}

