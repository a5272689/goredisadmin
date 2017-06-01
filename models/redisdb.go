package models

import (
	"crypto/md5"
	"strconv"
	"fmt"
	"encoding/json"
	"goredisadmin/utils"
	"strings"
)

func  GetHashName(hostname string,port int) (string) {
	h:=md5.New()
	h.Write([]byte(hostname))
	h.Write([]byte(strconv.Itoa(port)))
	hashname:=fmt.Sprintf("%x", h.Sum(nil))
	return hashname
}

type RedisInfo struct {
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Password string `json:"password"`
	Hashname string `json:"hashname"`
}

type RedissData struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	ConnectionStatus bool `json:"connection_status"`
	AuthStatus bool `json:"auth_status"`
	PingStatus bool `json:"ping_status"`
	Version string `json:"version"`
	Role string `json:"role"`
	UptimeInDays int `json:"uptime_in_days"`
	UsedMemoryRss int `json:"used_memory_rss"`
	Keys int `json:"keys"`
}

func GetRediss(redisinfos ...RedisInfo) []RedissData {
	Redis.Select(0)
	rediss:=[]RedissData{}
	newredisinfos:=[]RedisInfo{}
	if len(redisinfos)>0{
		for _,redisinfo:=range redisinfos{
			redisinfo.Hashname=GetHashName(redisinfo.Hostname,redisinfo.Port)
			exists,_:=Redis.Hexists("goredisadmin:rediss:hash",redisinfo.Hashname)
			if !exists {
				redisinfo.Save()
			}else {
				redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",redisinfo.Hashname)
				json.Unmarshal([]byte(redisinfostr),&redisinfo)
			}
			newredisinfos=append(newredisinfos,redisinfo)
		}
	}else {
		utils.Logger.Println("获取所有redis")
		redisslist,err:=Redis.Hkeys("goredisadmin:rediss:hash")
		if err!=nil{
			return rediss
		}
		for _,tmphashname:=range redisslist{
			redisinfo:=&RedisInfo{}
			redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",tmphashname)
			json.Unmarshal([]byte(redisinfostr),redisinfo)
			newredisinfos=append(newredisinfos,*redisinfo)
		}
	}
	for id,redisinfo:=range newredisinfos {
		redisinfojson,_:=json.Marshal(redisinfo)
		utils.Logger.Println(string(redisinfojson))
		redisC, err, conn, auth, ping := NewRedis(redisinfo.Hostname, redisinfo.Port, redisinfo.Password)
		var version,role string
		var uptime_in_days,used_memory_rss,keys int
		if err == nil {
			info, _ := redisC.Info()
			version = info["redis_version"]
			role = info["role"]
			uptime_in_days,_ = strconv.Atoi(info["uptime_in_days"])
			used_memory_rss,err = strconv.Atoi(info["used_memory_rss"])
			used_memory_rss=used_memory_rss/8/1024
			dbsinfo, _ := redisC.Info("Keyspace")
			for _,dbinfo:=range dbsinfo{
				keyinfolist:=strings.Split(dbinfo,",")
				infolist:=strings.Split(keyinfolist[0],"=")
				keysnum,_:=strconv.Atoi(infolist[1])
				keys+=keysnum
			}

		}
		rediss = append(rediss, RedissData{Id: id, Hostname: redisinfo.Hostname, Port: redisinfo.Port,UptimeInDays:uptime_in_days,
			ConnectionStatus:conn, AuthStatus: auth, PingStatus: ping, Version: version, Role: role,UsedMemoryRss:used_memory_rss,Keys:keys})
	}
	return rediss
}

func GetRedisNames() ([]string) {
	Redis.Select(0)
	redis_list:=[]string{}
	redisslist,err:=Redis.Hkeys("goredisadmin:rediss:hash")
	if err!=nil{
		return redis_list
	}
	for _,tmphashname:=range redisslist{
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		redis_list=append(redis_list,fmt.Sprintf("%v:%v",redisinfo.Hostname,redisinfo.Port))
	}
	return redis_list
}

func GetRedisDbs(rediss []string) (map[string][]string) {
	Redis.Select(0)
	redis_db_map:=make(map[string][]string)
	for _,redis:=range rediss{
		redislist:=strings.Split(redis,":")
		tmpport,_:=strconv.Atoi(redislist[1])
		tmphashname:=GetHashName(redislist[0],tmpport)
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		tmp_redis_obj,err,_,_,_:=NewRedis(redisinfo.Hostname,redisinfo.Port,redisinfo.Password)
		redis_db_map[redis]=[]string{}
		if err!=nil{
			continue
		}
		databases_str,_:=tmp_redis_obj.ConfigGet("databases")
		databases,_:=strconv.Atoi(databases_str["databases"])
		for dbnum:=0;dbnum<databases;dbnum++{
			redis_db_map[redis]=append(redis_db_map[redis],strconv.Itoa(dbnum))
		}
	}
	return redis_db_map
}

func (r *RedisInfo)Save() (result bool,err error) {
	Redis.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	jsonstr,err:=json.Marshal(r)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}


func  (r *RedisInfo)Del() (bool,error) {
	Redis.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	_,err:=Redis.Hdel("goredisadmin:rediss:hash",r.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

type KeysData struct {
	Key string `json:"key"`
	//Type string `json:"type"`
	//Ttl int `json:"ttl"`
}

func (r *RedisInfo) GetKeys(pattern string,dbindex int) ([]KeysData) {
	keys:=[]KeysData{}
	if len(pattern)==0{
		return keys
	}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, err, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	if err!=nil{
		return keys
	}
	keyslist,_:=redisC.Keys(pattern)
	for _,keyname:=range keyslist{
		//ttl,_:=redisC.Ttl(keyname)
		//typestr,_:=redisC.Type(keyname)
		keys=append(keys,KeysData{Key:keyname})
	}
	utils.Logger.Println("keys",keys)
	return keys
}

func (r *RedisInfo) DelKeys(keyslist []string,dbindex int) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Del(keyname)
		if err==nil{
			delkeyslist=append(delkeyslist,keyname)
		}
	}
	utils.Logger.Println("delkeyslist",delkeyslist)
	return delkeyslist
}

func (r *RedisInfo) ExpireKeys(keyslist []string,seconds int,dbindex int) ([]string) {
	var expire_key_list=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Expire(keyname,seconds)
		if err==nil{
			expire_key_list=append(expire_key_list,keyname)
		}
	}
	utils.Logger.Println("expire_key_list",expire_key_list)
	return expire_key_list
}

func (r *RedisInfo) PersistKeys(keyslist []string,dbindex int) ([]string) {
	var persist_key_list=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Persist(keyname)
		if err==nil{
			persist_key_list=append(persist_key_list,keyname)
		}
	}
	utils.Logger.Println("persist_key_list",persist_key_list)
	return persist_key_list
}

func (r *RedisInfo) SetKey(key,val interface{},dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Set(key,val)
}

func (r *RedisInfo) HsetKey(key,field string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Hset(key,field,value)
}

func (r *RedisInfo) LpushKey(key string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Lpush(key,value)
}

func (r *RedisInfo) LsetKey(key string,index int,value interface{},dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Lset(key,index,value)
}


func (r *RedisInfo) SaddKey(key string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Sadd(key,value)
}

func (r *RedisInfo) ZaddKey(key string,score int,value interface{},dbindex int) (int, error) {
	utils.Logger.Println(key,score,value,dbindex)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Zadd(key,score,value)
}

func (r *RedisInfo) GetKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Get(key)
}


func (r *RedisInfo) HmgetKey(key string,dbindex int) (map[string]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Hgetall(key)
}


func (r *RedisInfo) LrangeKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Lrange(key,0,-1)
}

func (r *RedisInfo) SmembersKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Smembers(key)
}

func (r *RedisInfo) TtlKey(key string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Ttl(key)
}

func (r *RedisInfo) TypeKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Type(key)
}


func (r *RedisInfo) ZrangeKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Zrange(key,0,-1,true)
}

func (r *RedisInfo) RenameKey(key,newkey string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Renamenx(key,newkey)
}

func (r *RedisInfo) DelStrValKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Set(key,"")
}

func (r *RedisInfo) DelHashValKey(key string,field string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Hdel(key,field)
}

func (r *RedisInfo) DelListValKey(key string,index int,dbindex int) (error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Ldel(key,index)
}

func (r *RedisInfo) DelSetValKey(key,member string,dbindex int) (int,error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Srem(key,member)
}

func (r *RedisInfo) DelZsetValKey(key,member string,dbindex int) (int,error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Select(dbindex)
	return redisC.Zrem(key,member)
}