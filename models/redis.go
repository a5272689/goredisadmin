package models

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"errors"
	"goredisadmin/utils"
)

func NewRedis() (*redis.Client,error)  {
	rc:=utils.Rc
	client, err := redis.Dial("tcp", fmt.Sprintf("%v:%v",rc.Host,rc.Port))
	if err!=nil{
		return client,err
	}
	if rc.Passwd!=""{
		rs:=client.Cmd("auth",rc.Passwd)
		result,_:=rs.Str()
		if result!="OK"{
			utils.Logger.Println("redis %v:%v 认证失败！！！",rc.Host,rc.Port)
		}
	}
	rs:=client.Cmd("ping")
	result,_:=rs.Str()
	if result!="PONG"{
		utils.Logger.Println("redis %v:%v 连接失败！！！，ping结果：%v",rc.Host,rc.Port,result)
	}
	return client,err
}

var Redis,_=NewRedis()

func CheckredisResult(result string,err error) (error) {
	if err!=nil{
		return err
	}
	if result!="ok"{
		return errors.New(result)
	}
	return nil
}