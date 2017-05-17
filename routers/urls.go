package routers

import (
	"github.com/gorilla/mux"
	"goredisadmin/controllers"
)


func Urls() *mux.Router  {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.MainHandler)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/login",controllers.Login)
	r.HandleFunc("/loginauth",controllers.LoginAuth)
	return r
}


