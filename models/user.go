package models

import (
	"crypto/sha256"
	"fmt"
)

type User struct {
	UserName string
	Role string
	PassWord string
}

func (u *User)GetPassWord() (dbpass string,err error)  {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	dbpass,err=Redis.Client.Get("goredisadmin:user:"+u.UserName)
	u.PassWord=dbpass
	return dbpass,err
}

func (u *User)GetRole() (role string,err error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	role,err=Redis.Client.Get("goredisadmin:userrole:"+u.UserName)
	u.Role=role
	return role,err
}

func  (u *User)HashPasswd(passwd string) (string) {
	h:=sha256.New()
	h.Write([]byte(passwd))
	h.Write([]byte(string(len(passwd))))
	h.Write([]byte("goredisadmin"))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func  (u *User)ChangePasswd(passwd string) (error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	hashpasswd:=u.HashPasswd(passwd)
	return CheckredisResult(Redis.Client.Set("goredisadmin:user:"+u.UserName,hashpasswd))
}

func  (u *User)ChangeRole(role string) (error) {
	Redis.Lock()
	defer Redis.Unlock()
	Redis.Client.Select(0)
	return CheckredisResult(Redis.Client.Set("goredisadmin:userrole:"+u.UserName,role))
}

