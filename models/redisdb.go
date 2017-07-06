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
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Password string `json:"password"`
	ConnectionStatus bool `json:"connection_status"`
	Version string `json:"version"`
	Role string `json:"role"`
	UptimeInDays int `json:"uptime_in_days"`
	UsedMemoryRss int `json:"used_memory_rss"`
	Keys int `json:"keys"`
	Dbs int `json:"dbs"`
	Group string `json:"group"`
	Mastername string `json:"mastername"`
	Hashname string `json:"hashname"`
	Remark string `json:"remark"`
}

func GetRediss(redisinfos ...RedisInfo) (rediss []RedisInfo) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return rediss
	}
	redisClient.Select(0)
	if len(redisinfos)>0{
		for _,redisinfo:=range redisinfos{
			redisinfo.Hashname=GetHashName(redisinfo.Hostname,redisinfo.Port)
			exists,_:=redisClient.Hexists("goredisadmin:rediss:hash",redisinfo.Hashname)
			if exists {
				tmpredisinfo:=&RedisInfo{}
				redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",redisinfo.Hashname)
				json.Unmarshal([]byte(redisinfostr),tmpredisinfo)
				redisinfo=*tmpredisinfo
			}
			rediss=append(rediss,redisinfo)
		}
	}else {
		utils.Logger.Println("获取所有redis")
		redisslist,err:=redisClient.Hkeys("goredisadmin:rediss:hash")
		if err!=nil{
			return rediss
		}
		for id,tmphashname:=range redisslist{
			redisinfo:=&RedisInfo{}
			redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",tmphashname)
			json.Unmarshal([]byte(redisinfostr),redisinfo)
			redisinfo.Id=id
			redisinfo.Password=""
			rediss=append(rediss,*redisinfo)
		}
	}
	return rediss
}

func GetRedisNames() (redis_list []string) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return redis_list
	}
	redisClient.Select(0)
	redisslist,err:=redisClient.Hkeys("goredisadmin:rediss:hash")
	if err!=nil{
		return redis_list
	}
	for _,tmphashname:=range redisslist{
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		redis_list=append(redis_list,fmt.Sprintf("%v:%v",redisinfo.Hostname,redisinfo.Port))
	}

	return redis_list
}

func GetRedisDbs(rediss []string) (map[string][]string) {
	redis_db_map:=make(map[string][]string)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return redis_db_map
	}
	redisClient.Select(0)
	for _,redis:=range rediss{
		redis_db_map[redis]=[]string{}
		redislist:=strings.Split(redis,":")
		tmpport,_:=strconv.Atoi(redislist[1])
		tmphashname:=GetHashName(redislist[0],tmpport)
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		for dbnum:=0;dbnum<redisinfo.Dbs;dbnum++{
			redis_db_map[redis]=append(redis_db_map[redis],strconv.Itoa(dbnum))
		}
	}
	utils.Logger.Println(redis_db_map)
	return redis_db_map
}

func (r *RedisInfo)Save() (result bool,err error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	jsonstr,err:=json.Marshal(r)
	if err!=nil{
		return false,err
	}
	_,err=redisClient.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func (r *RedisInfo)ChangePassword() (result bool,err error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	tmpRedisInfo:=&RedisInfo{}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),tmpRedisInfo)
	tmpRedisInfo.Password=r.Password
	utils.Logger.Println("new:",r,"now:",tmpRedisInfo)
	jsonstr,err:=json.Marshal(tmpRedisInfo)
	if err!=nil{
		return false,err
	}
	_,err=redisClient.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	go UpdateRedisInfo(tmpRedisInfo.Hostname,tmpRedisInfo.Port,tmpRedisInfo.Hashname,"")
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}


func (r *RedisInfo)Change() (result bool,err error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	tmpRedisInfo:=&RedisInfo{}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),tmpRedisInfo)
	tmpRedisInfo.Mastername=r.Mastername
	tmpRedisInfo.Group=r.Group
	tmpRedisInfo.Remark=r.Remark
	utils.Logger.Println("new:",r,"now:",tmpRedisInfo)
	jsonstr,err:=json.Marshal(tmpRedisInfo)
	if err!=nil{
		return false,err
	}
	_,err=redisClient.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (r *RedisInfo)Del() (bool,error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	_,err=redisClient.Hdel("goredisadmin:rediss:hash",r.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}
