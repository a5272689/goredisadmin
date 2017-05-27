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
	Redis.Select(0)
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
	Redis.Select(0)
	s.GetHashName()
	_,err:=Redis.Hdel("goredisadmin:sentinels:hash",s.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

type SentinelsData struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Masters []string `json:"masters"`
	ConnectionStatus bool `json:"connection_status"`
	MasterRediss map[string][]map[string]string `json:"master_rediss"`
	Version string `json:"version"`

}

func GetSentinels() []SentinelsData{
	Redis.Select(0)
	sentinels:=[]SentinelsData{}
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
		sentinel:=SentinelsData{Id:id,Hostname:tmpsentinel.Hostname,Port:tmpsentinel.Port,Version:version,
			Masters:masters,ConnectionStatus:ping,MasterRediss:masterrediss}
		sentinels=append(sentinels,sentinel)
	}
	return sentinels
}

