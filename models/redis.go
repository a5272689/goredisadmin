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
	defer client.Close()
	if err!=nil{
		return client,err
	}
	if rc.Passwd!=""{
		result,_:=client.Auth(rc.Passwd)
		if result!="OK"{
			utils.Logger.Println("redis %v:%v 认证失败！！！",rc.Host,rc.Port)
		}
	}
	result,err:=client.Ping()
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