package models


import (
	"strconv"
	"fmt"
	"encoding/json"
	"goredisadmin/utils"
	"strings"
)


type RoleInfo struct {
	Role string
	Slaves []*RedisInfo
}

func (r *RedisInfo) GetRoleInfo() (*RoleInfo) {
	roleinfo:=&RoleInfo{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return roleinfo
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return roleinfo
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		roleinfo.Role="slave"
		return roleinfo
	}
	ReplicationInfo,_:=redisC.Info("Replication")
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
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return keys
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return keys
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return keys
	}
	redisC.Select(dbindex)
	keyslist,_:=redisC.Keys(pattern)
	for _,keyname:=range keyslist{
		keys=append(keys,KeysData{Key:keyname})
	}
	utils.Logger.Println("keys",keys)
	return keys
}


func (r *RedisInfo) DelKeys(keyslist []string,dbindex int) ([]string) {
	var delkeyslist=[]string{}
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return delkeyslist
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return delkeyslist
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return delkeyslist
	}
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Del(keyname)
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
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return expire_key_list
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return expire_key_list
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return expire_key_list
	}
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Expire(keyname,seconds)
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
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return persist_key_list
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return persist_key_list
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return persist_key_list
	}
	redisC.Select(dbindex)
	for _,keyname:=range keyslist{
		_,err:=redisC.Persist(keyname)
		if err==nil{
			persist_key_list=append(persist_key_list,keyname)
		}
	}
	utils.Logger.Println("persist_key_list",persist_key_list)
	return persist_key_list
}



func (r *RedisInfo) SetKey(key,val interface{},dbindex int) (string, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return "",err
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return "",err
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return "",err
	}
	redisC.Select(dbindex)
	return redisC.Set(key,val)
}

func (r *RedisInfo) HsetKey(key,field string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return 0,err
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return 0,err
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return 0,err
	}
	redisC.Select(dbindex)
	return redisC.Hset(key,field,value)
}

func (r *RedisInfo) LpushKey(key string,value interface{},dbindex int) (int, error) {
	r.Hashname=GetHashName(r.Hostname,r.Port)
	redisClient,err:=RedisPool.Get()
	defer RedisPool.Put(redisClient)
	if err!=nil{
		return 0,err
	}
	redisinfostr,_:=redisClient.Hget("goredisadmin:rediss:hash",r.Hashname)
	json.Unmarshal([]byte(redisinfostr),r)
	utils.Logger.Println(r)
	redisPool, err := NewPool(r.Hostname,r.Port, r.Password)
	if err!=nil{
		return 0,err
	}
	redisC,err:=redisPool.Get()
	defer redisPool.Put(redisC)
	if err!=nil{
		return 0,err
	}
	redisC.Select(dbindex)
	return redisC.Lpush(key,value)
}

//func (r *RedisInfo) LsetKey(key string,index int,value interface{},dbindex int) (string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Lset(key,index,value)
//}
//
//
//func (r *RedisInfo) SaddKey(key string,value interface{},dbindex int) (int, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Sadd(key,value)
//}
//
//func (r *RedisInfo) ZaddKey(key string,score int,value interface{},dbindex int) (int, error) {
//	utils.Logger.Println(key,score,value,dbindex)
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Zadd(key,score,value)
//}
//
//func (r *RedisInfo) GetKey(key string,dbindex int) (string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Get(key)
//}
//
//
//func (r *RedisInfo) HmgetKey(key string,dbindex int) (map[string]string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Hgetall(key)
//}
//
//
//func (r *RedisInfo) LrangeKey(key string,dbindex int) ([]string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Lrange(key,0,-1)
//}
//
//func (r *RedisInfo) SmembersKey(key string,dbindex int) ([]string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Smembers(key)
//}
//
//func (r *RedisInfo) TtlKey(key string,dbindex int) (int, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Ttl(key)
//}
//
//func (r *RedisInfo) TypeKey(key string,dbindex int) (string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Type(key)
//}
//
//
//func (r *RedisInfo) ZrangeKey(key string,dbindex int) ([]string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Zrange(key,0,-1,true)
//}
//
//func (r *RedisInfo) RenameKey(key,newkey string,dbindex int) (int, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Renamenx(key,newkey)
//}
//
//func (r *RedisInfo) DelStrValKey(key string,dbindex int) (string, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Set(key,"")
//}
//
//func (r *RedisInfo) DelHashValKey(key string,field string,dbindex int) (int, error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Hdel(key,field)
//}
//
//func (r *RedisInfo) DelListValKey(key string,index int,dbindex int) (error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Ldel(key,index)
//}
//
//func (r *RedisInfo) DelSetValKey(key,member string,dbindex int) (int,error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Srem(key,member)
//}
//
//func (r *RedisInfo) DelZsetValKey(key,member string,dbindex int) (int,error) {
//	r.Hashname=GetHashName(r.Hostname,r.Port)
//	Redis.Lock()
//	defer Redis.Unlock()
//	redisinfostr,_:=Redis.Client.Hget("goredisadmin:rediss:hash",r.Hashname)
//	json.Unmarshal([]byte(redisinfostr),r)
//	utils.Logger.Println(r)
//	redisC, _, _, _, _ := NewRedis(r.Hostname,r.Port, r.Password)
//	redisC.Lock()
//	defer redisC.Unlock()
//	redisC.Client.Select(dbindex)
//	return redisC.Client.Zrem(key,member)
//}
