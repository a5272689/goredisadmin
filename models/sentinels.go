package models

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"encoding/json"
)

type Sentinel struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Masters []string `json:"masters"`
	ConnectionStatus bool `json:"connection_status"`
	MasterRediss map[string][]map[string]string `json:"master_rediss"`
	Version string `json:"version"`
	Hashname string `json:"hashname"`

}

func  (s *Sentinel)GetHashName() (string) {
	h:=md5.New()
	h.Write([]byte(s.Hostname))
	h.Write([]byte(strconv.Itoa(s.Port)))
	s.Hashname=fmt.Sprintf("%x", h.Sum(nil))
	return s.Hashname
}


func  (s *Sentinel)Create() (bool,error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	s.GetHashName()
	s.MasterRediss=map[string][]map[string]string{}
	s.Masters=[]string{}
	jsonstr,err:=json.Marshal(s)
	if err!=nil{
		return false,err
	}
	_,err=redisClient.Hset("goredisadmin:sentinels:hash",s.Hashname,string(jsonstr))
	go CheckHandle()
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (s *Sentinel)Del() (bool,error) {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return false,err
	}
	redisClient.Select(0)
	s.GetHashName()
	_,err=redisClient.Hdel("goredisadmin:sentinels:hash",s.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}


func GetSentinels() (sentinels []Sentinel,err error){
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return sentinels,err
	}
	redisClient.Select(0)
	sentinelslist,_:=redisClient.Hkeys("goredisadmin:sentinels:hash")
	for id,sentinelHashName:=range sentinelslist{
		sentinelinfo,_:=redisClient.Hget("goredisadmin:sentinels:hash",sentinelHashName)
		tmpsentinel:=&Sentinel{}
		json.Unmarshal([]byte(sentinelinfo),tmpsentinel)
		tmpsentinel.Id=id
		sentinels=append(sentinels,*tmpsentinel)
	}
	return sentinels,err
}


