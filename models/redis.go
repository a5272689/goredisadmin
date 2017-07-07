package models

import (
	"github.com/mediocregopher/radix.v2/pool"
	"fmt"
	"goredisadmin/utils"
	"strconv"
	"errors"
	"time"
	"encoding/json"
	"strings"
)

func init()  {
	RedisPool.MaxSize=100
}

var RedisPoolMap=map[string]*pool.Pool{}
var redisPasswordMap=map[string]string{}

type ConnStatus struct {
	RedisPool *pool.Pool
	Err error
}

func GetConnStatus(host string,port int,passwd string,channel  chan ConnStatus)  {
	sentinelC,err:=NewPool(host,port,passwd)
	connstatus:=ConnStatus{RedisPool:sentinelC,Err:err}
	channel <- connstatus
	close(channel)
}

func NewPool(host string,port int,passwd string)(*pool.Pool, error){
	redisPoolKeyStr:=host+strconv.Itoa(port)
	oldRpool:=RedisPoolMap[redisPoolKeyStr]
	if oldRpool!=nil&&redisPasswordMap[redisPoolKeyStr]==passwd{
		return oldRpool,nil
	}
	rpool,err:=pool.New("tcp", fmt.Sprintf("%v:%v",host,port),passwd,5)
	if err==nil{
		RedisPoolMap[redisPoolKeyStr]=rpool
		redisPasswordMap[redisPoolKeyStr]=passwd
	}
	return rpool,err
}


var RedisPool,_=NewPool(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)

func CheckredisResult(result string,err error) (error) {
	if err!=nil{
		return err
	}
	if result!="OK"{
		return errors.New(result)
	}
	return nil
}



func CheckRedis()  {
	utils.Logger.Println("开始检测！！")
	for {
		go CheckHandle()
		time.Sleep(time.Minute*10)
	}


}

func CheckHandle()  {
	sentinels:=[]Sentinel{}
	redisClient,_:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	redisClient.Select(0)
	channels:=[]chan ConnStatus{}
	sentinelslist,_:=redisClient.Hkeys("goredisadmin:sentinels:hash")
	for _,sentinelHashName:=range sentinelslist{
		sentinelinfo,_:=redisClient.Hget("goredisadmin:sentinels:hash",sentinelHashName)
		tmpsentinel:=&Sentinel{}
		json.Unmarshal([]byte(sentinelinfo),tmpsentinel)
		tmpchannel := make(chan ConnStatus)
		go GetConnStatus(tmpsentinel.Hostname,tmpsentinel.Port,"",tmpchannel)
		channels=append(channels,tmpchannel)
		sentinels=append(sentinels,*tmpsentinel)
	}
	sentinelsRedisInfoHashKey:=map[string]RedisInfo{}
	for i,tmpchannel:=range channels{
		masters:=[]string{}
		masterrediss:=make(map[string][]map[string]string)
		var version string
		connstatus := <-tmpchannel
		if connstatus.Err==nil{
			tmpredisClient,err:=connstatus.RedisPool.Get()
			defer connstatus.RedisPool.Put(tmpredisClient)
			if err!=nil{
				continue
			}
			mastersinfo,_:=tmpredisClient.Masters()
			for _,masterinfo:=range mastersinfo{
				masters=append(masters,masterinfo["name"])
				mastermaster:=map[string]string{"hostname":masterinfo["ip"],"port":masterinfo["port"]}
				portInt,_:=strconv.Atoi(masterinfo["port"])
				hashName:=GetHashName(masterinfo["ip"],portInt)
				sentinelsRedisInfoHashKey[hashName]=RedisInfo{Hostname:masterinfo["ip"],Port:portInt,Hashname:hashName,Mastername:masterinfo["name"]}
				//UpdateRedisInfo(masterinfo["ip"],portInt,hashName,masterinfo["name"])
				redissinfo:=[]map[string]string{mastermaster}
				slavesinfo,_:=tmpredisClient.Slaves(masterinfo["name"])
				for _,slaveinfo:=range slavesinfo{
					tmpinfo:=map[string]string{"hostname":slaveinfo["ip"],"port":slaveinfo["port"]}
					portInt,_:=strconv.Atoi(slaveinfo["port"])
					hashName:=GetHashName(slaveinfo["ip"],portInt)
					sentinelsRedisInfoHashKey[hashName]=RedisInfo{Hostname:slaveinfo["ip"],Port:portInt,Hashname:hashName,Mastername:slaveinfo["name"]}
					//UpdateRedisInfo(slaveinfo["ip"],portInt,hashName,masterinfo["name"])
					redissinfo=append(redissinfo,tmpinfo)
				}
				masterrediss[masterinfo["name"]]=redissinfo
			}
			info,_:=tmpredisClient.Info("Server")
			version=info["redis_version"]
			sentinels[i].ConnectionStatus=true
		}else {
			sentinels[i].ConnectionStatus=false
		}
		sentinels[i].Version=version
		sentinels[i].Masters=masters
		sentinels[i].MasterRediss=masterrediss
		jsonstr,err:=json.Marshal(sentinels[i])
		if err!=nil{
			continue
		}
		utils.Logger.Println(sentinels[i])
		redisClient.Hset("goredisadmin:sentinels:hash",sentinels[i].Hashname,string(jsonstr))
	}
	allRedisInfo,_:=redisClient.Hgetall("goredisadmin:rediss:hash")
	for redisInfoHashKey,redisInfo:=range sentinelsRedisInfoHashKey{
		delete(allRedisInfo,redisInfoHashKey)
		UpdateRedisInfo(redisInfo.Hostname,redisInfo.Port,redisInfo.Hashname,redisInfo.Mastername)
	}
	for hashName,redisInfoStr:=range allRedisInfo{
		redisinfo:=&RedisInfo{}
		json.Unmarshal([]byte(redisInfoStr),redisinfo)
		UpdateRedisInfo(redisinfo.Hostname,redisinfo.Port,hashName,redisinfo.Mastername)
	}
}

func UpdateRedisInfo(host string,port int,hashName,masterName string)  {
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return
	}
	redisClient.Select(0)
	exists,_:=redisClient.Hexists("goredisadmin:rediss:hash",hashName)
	redisinfo:=&RedisInfo{}
	if exists{
		redisinfoStr,_:=redisClient.Hget("goredisadmin:rediss:hash",hashName)
		json.Unmarshal([]byte(redisinfoStr),redisinfo)
	}else {
		redisinfo.Hostname=host
		redisinfo.Port=port
	}
	if len(masterName)>0{
		redisinfo.Mastername=masterName
	}
	redisPool,err:=NewPool(host,port,redisinfo.Password)
	if err!=nil{
		redisinfo.ConnectionStatus=false
	}else {
		redisinfo.ConnectionStatus=true
		tmpRedisClient,err:=redisPool.Get()
		defer redisPool.Put(tmpRedisClient)
		if err==nil{
			info, _ := tmpRedisClient.Info()
			redisinfo.Version= info["redis_version"]
			redisinfo.Role= info["role"]
			redisinfo.UptimeInDays,_ = strconv.Atoi(info["uptime_in_days"])
			used_memory_rss,_ := strconv.Atoi(info["used_memory_rss"])
			redisinfo.UsedMemoryRss=used_memory_rss/8/1024
			dbsinfo, _ := tmpRedisClient.Info("Keyspace")
			keys:=0
			for _,dbinfo:=range dbsinfo{
				keyinfolist:=strings.Split(dbinfo,",")
				infolist:=strings.Split(keyinfolist[0],"=")
				keysnum,_:=strconv.Atoi(infolist[1])
				keys+=keysnum
			}
			redisinfo.Keys=keys
			databases_str,_:=tmpRedisClient.ConfigGet("databases")
			redisinfo.Dbs,_=strconv.Atoi(databases_str["databases"])
		}
	}
	utils.Logger.Printf("主机：%v 端口：%v 密码：%v 更新",host,port,redisinfo.Password)
	redisinfo.Save()
}