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
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	Password string `json:"password"`
	Mastername string `json:"mastername"`
	Group string `json:"group"`
	Hashname string `json:"hashname"`
	Remark string `json:"remark"`
}

type RedissData struct {
	Id int `json:"id"`
	Hostname string `json:"hostname"`
	Port int `json:"port"`
	ConnectionStatus bool `json:"connection_status"`
	AuthStatus bool `json:"auth_status"`
	PingStatus bool `json:"ping_status"`
	Version string `json:"version"`
	Role string `json:"role"`
	UptimeInDays int `json:"uptime_in_days"`
	UsedMemoryRss int `json:"used_memory_rss"`
	Keys int `json:"keys"`
	Group string `json:"group"`
	Mastername string `json:"mastername"`
	Remark string `json:"remark"`
}

func GetRediss(redisinfos ...RedisInfo) []RedissData {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	rediss:=[]RedissData{}
	newredisinfos:=[]RedisInfo{}
	if len(redisinfos)>0{
		for _,redisinfo:=range redisinfos{
			redisinfo.Hashname=GetHashName(redisinfo.Hostname,redisinfo.Port)
			exists,_:=Redis.Client.Hexists("goredisadmin:rediss:hash",redisinfo.Hashname)
			if !exists {
				redisinfo.Save()
			}else {
				tmpredisinfo:=&RedisInfo{}
				redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",redisinfo.Hashname)
				json.Unmarshal([]byte(redisinfostr),tmpredisinfo)
				tmpredisinfo.Mastername=redisinfo.Mastername
				tmpredisinfo.Save()
				redisinfo=*tmpredisinfo
			}
			newredisinfos=append(newredisinfos,redisinfo)
		}
	}else {
		utils.Logger.Println("获取所有redis")
		redisslist,err:=Redis.Client.Hkeys("goredisadmin:rediss:hash")
		if err!=nil{
			return rediss
		}
		for _,tmphashname:=range redisslist{
			redisinfo:=&RedisInfo{}
			redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",tmphashname)
			json.Unmarshal([]byte(redisinfostr),redisinfo)
			newredisinfos=append(newredisinfos,*redisinfo)
		}
	}
	channels:=[]chan ConnStatus{}
	for id,redisinfo:=range newredisinfos {
		redisinfojson,_:=json.Marshal(redisinfo)
		utils.Logger.Println(string(redisinfojson))
		tmpchannel := make(chan ConnStatus)
		go GetConnStatus(redisinfo.Hostname, redisinfo.Port, redisinfo.Password,tmpchannel)
		channels=append(channels,tmpchannel)
		rediss = append(rediss, RedissData{Id: id, Hostname: redisinfo.Hostname, Port: redisinfo.Port,
			Mastername:redisinfo.Mastername,Group:redisinfo.Group,Remark:redisinfo.Remark})
	}
	for i,tmpchannel:=range channels{
		var version,role string
		var uptime_in_days,used_memory_rss,keys int
		connstatus := <-tmpchannel
		connstatus.Client.Lock()
		if connstatus.Err == nil {
			info, _ := connstatus.Client.Client.Info()
			version = info["redis_version"]
			role = info["role"]
			uptime_in_days,_ = strconv.Atoi(info["uptime_in_days"])
			used_memory_rss,_ = strconv.Atoi(info["used_memory_rss"])
			used_memory_rss=used_memory_rss/8/1024
			dbsinfo, _ := connstatus.Client.Client.Info("Keyspace")
			for _,dbinfo:=range dbsinfo{
				keyinfolist:=strings.Split(dbinfo,",")
				infolist:=strings.Split(keyinfolist[0],"=")
				keysnum,_:=strconv.Atoi(infolist[1])
				keys+=keysnum
			}

		}
		rediss[i].UptimeInDays=uptime_in_days
		rediss[i].ConnectionStatus=connstatus.Conn
		rediss[i].AuthStatus=connstatus.Auth
		rediss[i].PingStatus=connstatus.Ping
		rediss[i].Version=version
		rediss[i].Role=role
		rediss[i].UsedMemoryRss=used_memory_rss
		rediss[i].Keys=keys
		connstatus.Client.Unlock()
	}
	return rediss
}

func GetRedisNames() ([]string) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	redis_list:=[]string{}
	redisslist,err:=Redis.Client.Hkeys("goredisadmin:rediss:hash")
	if err!=nil{
		return redis_list
	}
	for _,tmphashname:=range redisslist{
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		redis_list=append(redis_list,fmt.Sprintf("%v:%v",redisinfo.Hostname,redisinfo.Port))
	}

	return redis_list
}

func GetRedisDbs(rediss []string) (map[string][]string) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	redis_db_map:=make(map[string][]string)
	channels:=[]chan ConnStatus{}
	redisinfos:=[]*RedisInfo{}
	for _,redis:=range rediss{
		redislist:=strings.Split(redis,":")
		tmpport,_:=strconv.Atoi(redislist[1])
		tmphashname:=GetHashName(redislist[0],tmpport)
		redisinfo:=&RedisInfo{}
		redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",tmphashname)
		json.Unmarshal([]byte(redisinfostr),redisinfo)
		redisinfos=append(redisinfos,redisinfo)
		tmpchannel := make(chan ConnStatus)

		go GetConnStatus(redisinfo.Hostname, redisinfo.Port, redisinfo.Password,tmpchannel)
		channels=append(channels,tmpchannel)
	}
	for i,redisinfo:=range redisinfos{
		redis:=fmt.Sprintf("%v:%v",redisinfo.Hostname,redisinfo.Port)
		connstatus := <-channels[i]
		connstatus.Client.Lock()
		defer connstatus.Client.Unlock()
		redis_db_map[redis]=[]string{}
		if connstatus.Err!=nil{
			continue
		}
		databases_str,_:=connstatus.Client.Client.ConfigGet("databases")
		databases,_:=strconv.Atoi(databases_str["databases"])
		for dbnum:=0;dbnum<databases;dbnum++{
			redis_db_map[redis]=append(redis_db_map[redis],strconv.Itoa(dbnum))
		}

	}
	return redis_db_map
}

func (r *RedisInfo)Save() (result bool,err error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	jsonstr,err:=json.Marshal(r)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Client.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func (r *RedisInfo)ChangePassword() (result bool,err error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	tmpRedisInfo:=&RedisInfo{}
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),tmpRedisInfo)
	tmpRedisInfo.Password=r.Password
	utils.Logger.Println("new:",r,"now:",tmpRedisInfo)
	jsonstr,err:=json.Marshal(tmpRedisInfo)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Client.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}


func (r *RedisInfo)Change() (result bool,err error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	tmpRedisInfo:=&RedisInfo{}
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),tmpRedisInfo)
	tmpRedisInfo.Mastername=r.Mastername
	tmpRedisInfo.Group=r.Group
	tmpRedisInfo.Remark=r.Remark
	utils.Logger.Println("new:",r,"now:",tmpRedisInfo)
	jsonstr,err:=json.Marshal(tmpRedisInfo)
	if err!=nil{
		return false,err
	}
	_,err=Redis.Client.Hset("goredisadmin:rediss:hash",r.Hashname,string(jsonstr))
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

func  (r *RedisInfo)Del() (bool,error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	_,err:=Redis.Client.Hdel("goredisadmin:rediss:hash",r.Hashname)
	if err!=nil{
		return false,err
	}else {
		return true,err
	}
}

type RoleInfo struct {
	Role string
	Slaves []*RedisInfo
}

func (r *RedisInfo) GetRoleInfo() (*RoleInfo) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, err, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	roleinfo:=&RoleInfo{}
	if err!=nil{
		roleinfo.Role="slave"
		return roleinfo
	}
	ReplicationInfo,_:=redisC.Client.Info("Replication")
	roleinfo.Role=ReplicationInfo["role"]
	if roleinfo.Role=="master"{
		connected_slaves,_:=strconv.Atoi(ReplicationInfo["connected_slaves"])

		for i:=0;i<connected_slaves;i++{
			tmpRedisInfo:=&RedisInfo{}
			slaveinfo:=ReplicationInfo[fmt.Sprintf("slave%v",i)]
			slaveinfoList:=strings.Split(slaveinfo,",")
			tmpRedisInfo.Hostname=strings.Split(slaveinfoList[0],"=")[1]
			tmpRedisInfo.Port,_=strconv.Atoi(strings.Split(slaveinfoList[1],"=")[1])
			roleinfo.Slaves=append(roleinfo.Slaves,tmpRedisInfo)
		}
	}
	return roleinfo
}


type KeysData struct {
	Key string `json:"key"`
	//Type string `json:"type"`
	//Ttl int `json:"ttl"`
}

func (r *RedisInfo) GetKeys(pattern string,dbindex int) ([]KeysData) {
	keys:=[]KeysData{}
	if len(pattern)==0{
		return keys
	}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, err, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	if err!=nil{
		return keys
	}
	redisC.Client.Select(dbindex)
	keyslist,_:=redisC.Client.Keys(pattern)
	for _,keyname:=range keyslist{
		//ttl,_:=redisC.Ttl(keyname)
		//typestr,_:=redisC.Type(keyname)
		keys=append(keys,KeysData{Key:keyname})
	}
	utils.Logger.Println("keys",keys)
	return keys
}

func (r *RedisInfo) DelKeys(keyslist []string,dbindex int) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Client.Del(keyname)
		if err==nil{
			delkeyslist=append(delkeyslist,keyname)
		}
	}
	utils.Logger.Println("delkeyslist",delkeyslist)
	return delkeyslist
}

func (r *RedisInfo) ExpireKeys(keyslist []string,seconds int,dbindex int) ([]string) {
	var expire_key_list=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Client.Expire(keyname,seconds)
		if err==nil{
			expire_key_list=append(expire_key_list,keyname)
		}
	}
	utils.Logger.Println("expire_key_list",expire_key_list)
	return expire_key_list
}

func (r *RedisInfo) PersistKeys(keyslist []string,dbindex int) ([]string) {
	var persist_key_list=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Client.Persist(keyname)
		if err==nil{
			persist_key_list=append(persist_key_list,keyname)
		}
	}
	utils.Logger.Println("persist_key_list",persist_key_list)
	return persist_key_list
}

func (r *RedisInfo) SetKey(key,val interface{},dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Set(key,val)
}

func (r *RedisInfo) HsetKey(key,field string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Hset(key,field,value)
}

func (r *RedisInfo) LpushKey(key string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Lpush(key,value)
}

func (r *RedisInfo) LsetKey(key string,index int,value interface{},dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Lset(key,index,value)
}


func (r *RedisInfo) SaddKey(key string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Sadd(key,value)
}

func (r *RedisInfo) ZaddKey(key string,score int,value interface{},dbindex int) (int, error) {
	utils.Logger.Println(key,score,value,dbindex)
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Zadd(key,score,value)
}

func (r *RedisInfo) GetKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Get(key)
}


func (r *RedisInfo) HmgetKey(key string,dbindex int) (map[string]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Hgetall(key)
}


func (r *RedisInfo) LrangeKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Lrange(key,0,-1)
}

func (r *RedisInfo) SmembersKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Smembers(key)
}

func (r *RedisInfo) TtlKey(key string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Ttl(key)
}

func (r *RedisInfo) TypeKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Type(key)
}


func (r *RedisInfo) ZrangeKey(key string,dbindex int) ([]string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Zrange(key,0,-1,true)
}

func (r *RedisInfo) RenameKey(key,newkey string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Renamenx(key,newkey)
}

func (r *RedisInfo) DelStrValKey(key string,dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Set(key,"")
}

func (r *RedisInfo) DelHashValKey(key string,field string,dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Hdel(key,field)
}

func (r *RedisInfo) DelListValKey(key string,index int,dbindex int) (error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Ldel(key,index)
}

func (r *RedisInfo) DelSetValKey(key,member string,dbindex int) (int,error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Srem(key,member)
}

func (r *RedisInfo) DelZsetValKey(key,member string,dbindex int) (int,error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	Redis.Lock()
	defer Redis.Unlock()
	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
	redisC.Lock()
	defer redisC.Unlock()
	redisC.Client.Select(dbindex)
	return redisC.Client.Zrem(key,member)
}