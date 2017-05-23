package models

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"errors"
	"strings"
)

type Sentinel struct {
	Sentinelid string
	HostName string
	Port int
	//Sentinel_cluster_name string
	HashName string
	//Masters []string
	//ConnectionStatus
}


func  (s *Sentinel)GetHashName() (string) {
	h:=md5.New()
	h.Write([]byte(s.HostName))
	h.Write([]byte(strconv.Itoa(s.Port)))
	s.HashName=fmt.Sprintf("%x", h.Sum(nil))
	return s.HashName
}

func (s *Sentinel)Save() (result bool,err error) {
	s.GetHashName()
	exists,_:=Redis.Hexists("goredisadmin:sentinels:hash",s.HashName+":hostname")
	if exists&&s.Sentinelid==""{
		return false,errors.New(fmt.Sprintf("hostname:%v,port:%v 已经存在无法新建!!!",s.HostName,s.Port))
	}
	if s.Sentinelid==""{
		return s.Create()
	}
	return s.Update()
}

func  (s *Sentinel)Create() (bool,error) {
	_,err:=Redis.Lpush("goredisadmin:sentinels:list",s.HashName)
	if err!=nil{
		return false,err
	}
	hmsetresult,err:=Redis.Hmset("goredisadmin:sentinels:hash",s.HashName+":hostname",s.HostName,s.HashName+":port",s.Port)
	if hmsetresult!="OK"||err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (s *Sentinel)Update() (bool,error) {
	sentinelid,err:=strconv.Atoi(s.Sentinelid)
	if err!=nil{
		return false,err
	}
	oldhashname,err:=Redis.Lindex("goredisadmin:sentinels:list",sentinelid)
	if oldhashname!=s.HashName{
		_,err:=Redis.Lset("goredisadmin:sentinels:list",sentinelid,s.HashName)
		if err!=nil{
			return false,err
		}
		_,err=Redis.Hdel("goredisadmin:sentinels:hash",oldhashname+":hostname",oldhashname+":port")
		fmt.Println(err)
		if err!=nil {
			return false,err
		}
	}
	hmsetresult,err:=Redis.Hmset("goredisadmin:sentinels:hash",s.HashName+":hostname",s.HostName,s.HashName+":port",s.Port)
	if hmsetresult!="OK"||err!=nil{
		return false,err
	}else {
		return true,err
	}


}

func GetSentinels() []map[string]interface{} {

	sentinels:=[]map[string]interface{}{}
	llen,err:=Redis.Llen("goredisadmin:sentinels:list")
	if err!=nil{
		return sentinels
	}
	sentinelslist,err:=Redis.Lrange("goredisadmin:sentinels:list",0,llen)
	if err!=nil{
		return sentinels
	}
	for id,sentinelHashName:=range sentinelslist{
		sentinelinfo,err:=Redis.Hmget("goredisadmin:sentinels:hash",sentinelHashName+":hostname",sentinelHashName+":port")
		if err!=nil{
			continue
		}
		port,_:=strconv.Atoi(sentinelinfo[1])
		sentinelC,err,_,_,ping:=NewRedis(sentinelinfo[0],port,"")
		masters:=[]string{}
		masterrediss:=make(map[string][]map[string]string)
		var version string
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
			infos:=strings.Split(info,"\n")
			for _,infostr:=range infos{
				infolist:=strings.Split(infostr,":")
				if len(infolist)==2{
					if infolist[0]=="redis_version"{
						version=infolist[1]
					}
				}
			}

		}
		sentinel:=map[string]interface{}{"id":id,"hostname":sentinelinfo[0],"port":port,"version":version,
			"masters":masters,"connection_status":ping,"master_rediss":masterrediss}
		sentinels=append(sentinels,sentinel)
	}
	return sentinels
}

