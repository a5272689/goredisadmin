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
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	s.GetHashName()
	jsonstr,err:=json.Marshal(s)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Client.Hset("goredisadmin:sentinels:hash",s.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (s *Sentinel)Del() (bool,error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	s.GetHashName()
	_,err:=Redis.Client.Hdel("goredisadmin:sentinels:hash",s.Hashname)
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
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	sentinels:=[]SentinelsData{}
	channels:=[]chan ConnStatus{}
	sentinelslist,err:=Redis.Client.Hkeys("goredisadmin:sentinels:hash")
	if err!=nil{
		return sentinels
	}

	for id,sentinelHashName:=range sentinelslist{
		sentinelinfo,_:=Redis.Client.Hget("goredisadmin:sentinels:hash",sentinelHashName)
		tmpsentinel:=&Sentinel{}
		json.Unmarshal([]byte(sentinelinfo),tmpsentinel)
		tmpchannel := make(chan ConnStatus)
		go GetConnStatus(tmpsentinel.Hostname,tmpsentinel.Port,"",tmpchannel)
		channels=append(channels,tmpchannel)
		sentinel:=SentinelsData{Id:id,Hostname:tmpsentinel.Hostname,Port:tmpsentinel.Port}
		sentinels=append(sentinels,sentinel)
	}
	for i,tmpchannel:=range channels{
		masters:=[]string{}
		masterrediss:=make(map[string][]map[string]string)
		var version string
		connstatus := <-tmpchannel
		connstatus.Client.Lock()
		defer connstatus.Client.Unlock()
		if connstatus.Err==nil{
			mastersinfo,_:=connstatus.Client.Client.Masters()
			for _,masterinfo:=range mastersinfo{
				masters=append(masters,masterinfo["name"])
				mastermaster:=map[string]string{"hostname":masterinfo["ip"],"port":masterinfo["port"]}
				redissinfo:=[]map[string]string{mastermaster}
				slavesinfo,_:=connstatus.Client.Client.Slaves(masterinfo["name"])
				for _,slaveinfo:=range slavesinfo{
					tmpinfo:=map[string]string{"hostname":slaveinfo["ip"],"port":slaveinfo["port"]}
					redissinfo=append(redissinfo,tmpinfo)
				}
				masterrediss[masterinfo["name"]]=redissinfo
			}
			info,_:=connstatus.Client.Client.Info("Server")
			version=info["redis_version"]
		}
		sentinels[i].Version=version
		sentinels[i].Masters=masters
		sentinels[i].ConnectionStatus=connstatus.Ping
		sentinels[i].MasterRediss=masterrediss
	}

	return sentinels
}


