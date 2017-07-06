package controllers



import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	//"fmt"
	//"time"
	//"encoding/json"
	"github.com/flosch/pongo2"
	"goredisadmin/models"
	//"goredisadmin/utils"
	//"github.com/bitly/go-simplejson"
	//"strconv"
	//"strings"
	//"io/ioutil"
	//"goredisadmin/utils"
	"fmt"
)

func initconText(r *http.Request) (pongo2.Context) {
	session := sessions.GetSession(r)
	redisClient,err:=models.RedisPool.Get()
	defer models.RedisPool.Put(redisClient)
	if err!=nil{
		return pongo2.Context{}
	}
	sentinels_keys,_:=redisClient.Hkeys("goredisadmin:sentinels:hash")
	redis_keys,_:=redisClient.Hkeys("goredisadmin:rediss:hash")
	user:=session.Get("user")
	username:=session.Get("username")
	if username==nil{
		username=user
	}
	userrole:=session.Get("role")
	if userrole==nil{
		userrole=""
	}
	return pongo2.Context{"username":username,"userrole":userrole,"urlpath":r.URL.Path,"sentinels":len(sentinels_keys),"redis":len(redis_keys)}
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

//func KeyValDelAPI(w http.ResponseWriter, r *http.Request) {
//	data, _ := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//	session := sessions.GetSession(r)
//	utils.Logger.Printf("[info] KeyValDelAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
//	jsonob,_:=simplejson.NewJson(data)
//	key,_:=jsonob.Get("key").String()
//	key_type,_:=jsonob.Get("type").String()
//	redisstr,_:=jsonob.Get("redis").String()
//	field,_:=jsonob.Get("field").String()
//	member,_:=jsonob.Get("val").String()
//	indexstr,_:=jsonob.Get("index").String()
//	index,_:=strconv.Atoi(indexstr)
//	redislist:=strings.Split(redisstr,":")
//	redisport,_:=strconv.Atoi(redislist[1])
//	redis_db,_:=jsonob.Get("redis_db").String()
//	redis_db_index,_:=strconv.Atoi(redis_db)
//	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
//	result:=new(JsonResult)
//	result.Result=true
//	switch key_type {
//	case "string":
//		_,err:=redisinfo.DelStrValKey(key,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "hash":
//		_,err:=redisinfo.DelHashValKey(key,field,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "list":
//		err:=redisinfo.DelListValKey(key,index,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "set":
//		_,err:=redisinfo.DelSetValKey(key,member,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "zset":
//		_,err:=redisinfo.DelZsetValKey(key,member,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	default:
//		result.Result=false
//		result.Info="不支持的类型"
//	}
//	utils.Logger.Println(key,key_type,redisport,redis_db_index)
//	jsonresult,_:=json.Marshal(result)
//	fmt.Fprint(w,string(jsonresult))
//}
//
//func KeySaveAPI(w http.ResponseWriter, r *http.Request) {
//	data, _ := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//	session := sessions.GetSession(r)
//	utils.Logger.Printf("[info] KeySaveAPI收到json串：%v,操作用户:%v",string(data),session.Get("user"))
//	jsonob,_:=simplejson.NewJson(data)
//	key_type,_:=jsonob.Get("type").String()
//	key,_:=jsonob.Get("key").String()
//	val,_:=jsonob.Get("val").String()
//	field,_:=jsonob.Get("field").String()
//	indexstr,_:=jsonob.Get("index").String()
//	scorestr,_:=jsonob.Get("score").String()
//	score,_:=strconv.Atoi(scorestr)
//	redisstr,_:=jsonob.Get("redis").String()
//	redislist:=strings.Split(redisstr,":")
//	redisport,_:=strconv.Atoi(redislist[1])
//	redis_db,_:=jsonob.Get("redis_db").String()
//	redis_db_index,_:=strconv.Atoi(redis_db)
//	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
//	result:=new(JsonResult)
//	result.Result=true
//	switch key_type {
//	case "string":
//		_,err:=redisinfo.SetKey(key,val,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "hash":
//		_,err:=redisinfo.HsetKey(key,field,val,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "list":
//		var err error
//		if indexstr==""{
//			_,err=redisinfo.LpushKey(key,val,redis_db_index)
//		}else {
//			index,_:=strconv.Atoi(indexstr)
//			_,err=redisinfo.LsetKey(key,index,val,redis_db_index)
//		}
//
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "set":
//		_,err:=redisinfo.SaddKey(key,val,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	case "zset":
//		_,err:=redisinfo.ZaddKey(key,score,val,redis_db_index)
//		if err!=nil{
//			result.Result=false
//			result.Info=fmt.Sprintf("报错：%v",err)
//		}
//	default:
//		result.Result=false
//		result.Info="不支持的类型"
//	}
//	jsonresult,_:=json.Marshal(result)
//	fmt.Fprint(w,string(jsonresult))
//}
//
//
//
//func KeyDataAPI(w http.ResponseWriter, r *http.Request) {
//	data, _ := ioutil.ReadAll(r.Body)
//	defer r.Body.Close()
//	utils.Logger.Println(string(data))
//	jsonob,_:=simplejson.NewJson(data)
//	keystr,_:=jsonob.Get("key").String()
//	redisstr,_:=jsonob.Get("redis").String()
//	redislist:=strings.Split(redisstr,":")
//	redisport,_:=strconv.Atoi(redislist[1])
//	redis_db,_:=jsonob.Get("redis_db").String()
//	redis_db_index,_:=strconv.Atoi(redis_db)
//	redisinfo:=models.RedisInfo{Hostname:redislist[0],Port:redisport}
//	key_type,err:=redisinfo.TypeKey(keystr,redis_db_index)
//	utils.Logger.Println(key_type,err)
//	key_ttl,err:=redisinfo.TtlKey(keystr,redis_db_index)
//	utils.Logger.Println(key_ttl,err)
//	alldata:=make(map[string]interface{})
//	rows:=[]map[string]string{}
//	alldata["type"]=key_type
//	alldata["ttl"]=key_ttl
//	switch key_type {
//	case "string":
//		val,err:=redisinfo.GetKey(keystr,redis_db_index)
//		utils.Logger.Println(val)
//		valmap:=make(map[string]string)
//		if err==nil{
//			valmap["val"]=val
//			rows=append(rows,valmap)
//
//		}
//	case "hash":
//		vals,err:=redisinfo.HmgetKey(keystr,redis_db_index)
//		if err==nil{
//			for field,val:=range vals{
//				valmap:=make(map[string]string)
//				valmap["val"]=val
//				valmap["field"]=field
//				rows=append(rows,valmap)
//			}
//		}
//		utils.Logger.Println(vals)
//	case "list":
//		vals,err:=redisinfo.LrangeKey(keystr,redis_db_index)
//		if err==nil{
//			for index,val:=range vals {
//				valmap := make(map[string]string)
//				valmap["val"] = val
//				valmap["index"]=strconv.Itoa(index)
//				rows = append(rows, valmap)
//			}
//		}
//		utils.Logger.Println(vals)
//	case "set":
//		vals,err:=redisinfo.SmembersKey(keystr,redis_db_index)
//		if err==nil{
//			for index,val:=range vals {
//				valmap := make(map[string]string)
//				valmap["val"] = val
//				valmap["index"]=strconv.Itoa(index)
//				rows = append(rows, valmap)
//			}
//		}
//		utils.Logger.Println(vals)
//	case "zset":
//		vals,err:=redisinfo.ZrangeKey(keystr,redis_db_index)
//		if err==nil{
//			for index,val:=range vals {
//				if index%2==0{
//					valmap := make(map[string]string)
//					valmap["val"] = val
//					valmap["score"]=vals[index+1]
//					rows = append(rows, valmap)
//				}
//			}
//		}
//	}
//	alldata["rows"]=rows
//	jsonresult,_:=json.Marshal(alldata)
//	fmt.Fprint(w,string(jsonresult))
//}
//
