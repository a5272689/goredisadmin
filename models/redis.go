package models

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"errors"
	"goredisadmin/utils"
)

func NewRedis(host string,port int,passwd string) (client *redis.Client,err error,conn,auth,ping bool)  {
	client, err = redis.Dial("tcp", fmt.Sprintf("%v:%v",host,port))
	if err!=nil{
		return client,err,conn,auth,ping
	}
	conn=true
	if passwd!=""{
		result,_:=client.Auth(passwd)
		if result!="OK"{
			utils.Logger.Println("redis %v:%v 认证失败！！！",host,port)
			return client,err,conn,auth,ping
		}
	}
	auth=true
	result,err:=client.Ping()
	if result!="PONG"{
		utils.Logger.Println("redis %v:%v 连接失败！！！，ping结果：%v",host,port,result)
		return client,err,conn,auth,ping
	}
	ping=true
	return client,err,conn,auth,ping
}

var Redis,_,_,_,_=NewRedis(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)

func CheckredisResult(result string,err error) (error) {
	if err!=nil{
		return err
	}
	if result!="ok"{
		return errors.New(result)
	}
	return nil
}