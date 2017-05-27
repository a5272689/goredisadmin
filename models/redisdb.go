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


func (r *RedisInfo)Save() (result bool,err error) {
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
	Type string `json:"type"`
	Ttl int `json:"ttl"`
}

func (r *RedisInfo) GetKeys(pattern string) ([]KeysData) {
	keys:=[]KeysData{}
	if len(pattern)==0{
		return keys
	}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, err, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return keys
	}
	keyslist,_:=redisC.Keys(pattern)
	for _,keyname:=range keyslist{
		ttl,_:=redisC.Ttl(keyname)
		typestr,_:=redisC.Type(keyname)
		keys=append(keys,KeysData{Key:keyname,Ttl:ttl,Type:typestr})
	}
	return keys
}

func (r *RedisInfo) DelKeys(keyslist []string) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	for _,keyname:=range keyslist{
		_,err:=redisC.Del(keyname)
		if err==nil{
			delkeyslist=append(delkeyslist,keyname)
		}
	}
	return delkeyslist
}

func (r *RedisInfo) ExpireKeys(keyslist []string,seconds int) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	for _,keyname:=range keyslist{
		_,err:=redisC.Expire(keyname,seconds)
		if err==nil{
			delkeyslist=append(delkeyslist,keyname)
		}
	}
	return delkeyslist
}

func (r *RedisInfo) PersistKeys(keyslist []string) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	for _,keyname:=range keyslist{
		_,err:=redisC.Persist(keyname)
		if err==nil{
			delkeyslist=append(delkeyslist,keyname)
		}
	}
	return delkeyslist
}