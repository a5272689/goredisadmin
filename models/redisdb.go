package models

import (
	"crypto/md5"
	"strconv"
	"fmt"
	"encoding/json"
	"goredisadmin/utils"
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

func GetRediss(redisinfos ...RedisInfo) []map[string]interface{} {
	rediss:=[]map[string]interface{}{}
	newredisinfos:=[]RedisInfo{}
	if len(redisinfos)>0{
		for _,redisinfo:=range redisinfos{
			redisinfo.Hashname=GetHashName(redisinfo.Hostname,redisinfo.Port)
			exists,_:=Redis.Hexists("goredisadmin:rediss:hash",redisinfo.Hashname)
			if !exists {
				redisinfo.Save()
			}else {
				redisinfostr,_:=Redis.Hget("goredisadmin:rediss:hash",redisinfo.Hashname)
				json.Unmarshal([]byte(redisinfostr),redisinfo)
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
		redisC, err, conn, auth, ping := NewRedis(redisinfo.Hostname, redisinfo.Port, redisinfo.Password)
		version := ""
		role := ""
		if err == nil {
			info, _ := redisC.Info()
			version = info["redis_version"]
			role = info["role"]
		}
		rediss = append(rediss, map[string]interface{}{"id": id, "hostname": redisinfo.Hostname, "port": redisinfo.Port,
			"connection_status":                         conn, "auth_status": auth, "ping_status": ping, "version": version, "role": role})
	}
	return rediss
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

//func (r *RedisInfo)Update() (result bool,index int,err error) {
//	sentinelid,err:=strconv.Atoi(r.Id)
//	if err!=nil{
//		return false,0,err
//	}
//	oldhashname,err:=Redis.Lindex("goredisadmin:sentinels:list",sentinelid)
//	if oldhashname!=r.HashName{
//		_,err:=Redis.Lset("goredisadmin:sentinels:list",sentinelid,s.HashName)
//		if err!=nil{
//			return false,err
//		}
//		_,err=Redis.Hdel("goredisadmin:sentinels:hash",oldhashname+":hostname",oldhashname+":port")
//		fmt.Println(err)
//		if err!=nil {
//			return false,err
//		}
//	}
//	hmsetresult,err:=Redis.Hmset("goredisadmin:sentinels:hash",s.HashName+":hostname",s.HostName,s.HashName+":port",s.Port)
//	if hmsetresult!="OK"||err!=nil{
//		return false,err
//	}else {
//		return true,err
//	}
//}


func GetAllRediss() []map[string]interface{} {
	rediss:=[]map[string]interface{}{}
	llen,err:=Redis.Llen("goredisadmin:rediss:list")
	if err!=nil{
		return rediss
	}
	redisslist,err:=Redis.Lrange("goredisadmin:rediss:list",0,llen)
	if err!=nil{
		return rediss
	}
	for _,redisname:=range redisslist{
		fmt.Println(redisname)
	}
	return rediss
}
