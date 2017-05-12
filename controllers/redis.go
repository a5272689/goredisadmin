package controllers

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
)

func NewRedis() (*redis.Client,error)  {
	rc:=Rc
	client, err := redis.Dial("tcp", fmt.Sprintf("%v:%v",rc.Host,rc.Port))
	if err!=nil{
		return client,err
	}
	if rc.Passwd!=""{
		rs:=client.Cmd("auth",rc.Passwd)
		result,_:=rs.Str()
		if result!="OK"{
			Logger.Println("redis %v:%v 认证失败！！！",rc.Host,rc.Port)
		}
	}
	rs:=client.Cmd("ping")
	result,_:=rs.Str()
	if result!="PONG"{
		Logger.Println("redis %v:%v 连接失败！！！，ping结果：%v",rc.Host,rc.Port,result)
	}
	return client,err
}

var Redis,_=NewRedis()