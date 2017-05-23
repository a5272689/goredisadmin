package models

import (

)
import (
	"crypto/md5"
	"strconv"
	"fmt"
)

func  GetHashName(hostname string,port int) (string) {
	h:=md5.New()
	h.Write([]byte(hostname))
	h.Write([]byte(strconv.Itoa(port)))
	hashname:=fmt.Sprintf("%x", h.Sum(nil))
	return hashname
}

type RedisInfo struct {
	id int
	HostName string
	Port int
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
	redissmap:=map[string]map[string]interface{}{}
	for redisid,redisname:=range redisslist{
		redissmap[redisname]=map[string]interface{}{"id":redisid}
	}
	if len(redisinfos)>0{
		for _,redisname:=range redisslist{
			fmt.Println(redisname)
		}
	}
	return rediss
}

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
