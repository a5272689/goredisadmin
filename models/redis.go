package models

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"errors"
	"goredisadmin/utils"
	"strconv"
	"time"
	"sync"
)

type ConnStatus struct {
	Client *redisClient
	Err error
	Conn bool
	Auth bool
	Ping bool
}

func GetConnStatus(host string,port int,passwd string,channel  chan ConnStatus)  {
	sentinelC,err,conn,auth,ping:=NewRedis(host,port,passwd)
	connstatus:=ConnStatus{Client:sentinelC,Err:err,Conn:conn,Auth:auth,Ping:ping}
	channel <- connstatus
	close(channel)
}

func NewRedis(host string,port int,passwd string) (redisclient *redisClient,err error,conn,auth,ping bool)  {
	portstr:=strconv.Itoa(port)
	redisclient=RedisMap[host+portstr]
	if redisclient.Client!=nil{
		redisclient.Lock()
		result,err:=redisclient.Client.Ping()
		if result=="PONG"{
			return redisclient,err,true,true,true
		}
		redisclient.Unlock()
	}

	tmpredisclient:=&redisClient{}
	tmpredisclient.Client, err = redis.DialTimeout("tcp", fmt.Sprintf("%v:%v",host,port),time.Second*2)
	//client, err = redis.Dial("tcp", fmt.Sprintf("%v:%v",host,port))
	if err!=nil{
		utils.Logger.Printf("redis %v:%v 连接失败！！！",host,port)
		return tmpredisclient,err,conn,auth,ping
	}
	conn=true
	if passwd!=""{
		result,_:=tmpredisclient.Client.Auth(passwd)
		if result!="OK"{
			utils.Logger.Printf("redis %v:%v 认证失败！！！",host,port)
			return tmpredisclient,err,conn,auth,ping
		}
	}
	auth=true
	result,err:=tmpredisclient.Client.Ping()
	if result!="PONG"{
		utils.Logger.Printf("redis %v:%v ping失败！！！",host,port)
		if passwd==""{
			auth=false
		}
		return tmpredisclient,err,conn,auth,ping
	}
	ping=true
	RedisMap[host+portstr]=redisclient
	return tmpredisclient,err,conn,auth,ping
}

type redisClient struct {
	Client *redis.Client
	sync.Locker
}

var Redis,_,_,_,_=NewRedis(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)

var RedisMap=map[string]*redisClient{}

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
	defer  func() {
		if x := recover(); x != nil {
			utils.Logger.Println("直连redis连接失败，重新连接！！！")
			Redis,_,_,_,_=NewRedis(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)
			CheckRedis()
		}
	}()
	utils.Logger.Println("开始检测！！")
	for i:=1;i==1;{
		time.Sleep(time.Second*1)
		Redis.Lock()
		_,err:=Redis.Client.Ping()
		Redis.Unlock()
		if err!=nil{
			utils.Logger.Println("直连redis连接失败，重新连接！！！")
			Redis,_,_,_,_=NewRedis(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)
		}
	}


}