package models

import (

)
import (
	"crypto/md5"
	"strconv"
	"fmt"
	"errors"
)

func  GetHashName(hostname string,port int) (string) {
	h:=md5.New()
	h.Write([]byte(hostname))
	h.Write([]byte(strconv.Itoa(port)))
	hashname:=fmt.Sprintf("%x", h.Sum(nil))
	return hashname
}

type RedisInfo struct {
	Id string
	HostName string
	Port int
	PassWord string
	//Sentinel_cluster_name string
	HashName string
	//Masters []string
	//ConnectionStatus
}

func GetRediss(redisinfos ...RedisInfo) []map[string]interface{} {
	rediss:=[]map[string]interface{}{}
	llen,err:=Redis.Llen("goredisadmin:rediss:list")
	if err!=nil{
		return rediss
	}
	redisslist,err:=Redis.Lrange("goredisadmin:rediss:list",0,llen)
	if err!=nil{
		return rediss
	}
	fmt.Println(redisinfos)
	redissmap:=map[string]int{}
	for redisid,redisname:=range redisslist{
		redissmap[redisname]=redisid
	}
	//if len(redisinfos)>0{
	//	for _,redisinfo:=range redisinfos{
	//		redisinfo.HashName=GetHashName(redisinfo.HostName,redisinfo.Port)
	//		exists,_:=Redis.Hexists("goredisadmin:rediss:hash",redisinfo.HashName+":hostname")
	//		if exists{
	//			redisinfo.Id=strconv.Itoa(redissmap[redisinfo.HashName])
	//		}else {
	//			redisslist,err:=Redis.Lpush("goredisadmin:rediss:list",redisinfo.HashName)
	//		}
	//		fmt.Println(redisinfo.HashName,redissmap[redisinfo.HashName])
	//		rediss=append(rediss,map[string]interface{}{"id":redisinfo.Id})
	//	}
	//}

	fmt.Println(redissmap)
	fmt.Println(rediss)
	return rediss
}


func (r *RedisInfo)Save() (result bool,index int,err error) {
	r.HashName=GetHashName(r.HostName,r.Port)
	exists,_:=Redis.Hexists("goredisadmin:rediss:hash",r.HashName+":hostname")
	if exists&&r.Id==""{
		return false,0,errors.New(fmt.Sprintf("hostname:%v,port:%v 已经存在无法新建!!!",r.HostName,r.Port))
	}
	if r.Id==""{
		return r.Create()
	}
	return false,0,errors.New(fmt.Sprintf("hostname:%v,port:%v 已经存在无法新建!!!",r.HostName,r.Port))
	//return r.Update()
}

func (r *RedisInfo)Create() (result bool,index int,err error) {
	index,err=Redis.Lpush("goredisadmin:rediss:list",r.HashName)
	if err!=nil{
		return false,0,err
	}
	hmsetresult,err:=Redis.Hmset("goredisadmin:sentinels:hash",r.HashName+":hostname",r.HostName,r.HashName+":port",r.Port,r.HashName+":password",r.PassWord)
	if hmsetresult!="OK"||err!=nil{
		return false,index,err
	}else {
		return true,index,err
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
