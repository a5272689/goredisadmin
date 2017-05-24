package models

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"encoding/json"
)

type Sentinel struct {
	Hostname string `json:"hostname"`
	Port int `json:"port"`
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
	s.GetHashName()
	jsonstr,err:=json.Marshal(s)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Hset("goredisadmin:sentinels:hash",s.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (s *Sentinel)Del() (bool,error) {
	s.GetHashName()
	_,err:=Redis.Hdel("goredisadmin:sentinels:hash",s.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func GetSentinels() []map[string]interface{} {

	sentinels:=[]map[string]interface{}{}
	sentinelslist,err:=Redis.Hkeys("goredisadmin:sentinels:hash")
	if err!=nil{
		return sentinels
	}
	for id,sentinelHashName:=range sentinelslist{
		sentinelinfo,_:=Redis.Hget("goredisadmin:sentinels:hash",sentinelHashName)
		tmpsentinel:=&Sentinel{}
		masters:=[]string{}
		masterrediss:=make(map[string][]map[string]string)
		var version string
		json.Unmarshal([]byte(sentinelinfo),tmpsentinel)
		sentinelC,err,_,_,ping:=NewRedis(tmpsentinel.Hostname,tmpsentinel.Port,"")
		if err==nil{
			mastersinfo,_:=sentinelC.Masters()
			for _,masterinfo:=range mastersinfo{
				masters=append(masters,masterinfo["name"])
				mastermaster:=map[string]string{"hostname":masterinfo["ip"],"port":masterinfo["port"]}
				redissinfo:=[]map[string]string{mastermaster}
				slavesinfo,_:=sentinelC.Slaves(masterinfo["name"])
				for _,slaveinfo:=range slavesinfo{
					tmpinfo:=map[string]string{"hostname":slaveinfo["ip"],"port":slaveinfo["port"]}
					redissinfo=append(redissinfo,tmpinfo)
				}
				masterrediss[masterinfo["name"]]=redissinfo
			}
			info,_:=sentinelC.Info("Server")
			version=info["redis_version"]
		}
		sentinel:=map[string]interface{}{"id":id,"hostname":tmpsentinel.Hostname,"port":tmpsentinel.Port,"version":version,
			"masters":masters,"connection_status":ping,"master_rediss":masterrediss}
		sentinels=append(sentinels,sentinel)
	}
	return sentinels
}

