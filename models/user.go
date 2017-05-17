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
	dbpass,err=Redis.Cmd("get","goredisadmin:user:"+u.UserName).Str()
	u.PassWord=dbpass
	return dbpass,err
}

func (u *User)GetRole() (role string,err error) {
	role,err=Redis.Cmd("get","goredisadmin:userrole:"+u.UserName).Str()
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
	hashpasswd:=u.HashPasswd(passwd)
	return CheckredisResult(Redis.Cmd("set","goredisadmin:user:"+u.UserName,hashpasswd).Str())
}

func  (u *User)ChangeRole(role string) (error) {
	return CheckredisResult(Redis.Cmd("set","goredisadmin:userrole:"+u.UserName,role).Str())
}

