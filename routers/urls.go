package routers

import (
	"github.com/gorilla/mux"
	"goredisadmin/controllers"
)


func Urls() *mux.Router  {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.MainHandler)
	r.HandleFunc("/a", controllers.Main2Handler)
	r.HandleFunc("/login",controllers.Login)
	r.HandleFunc("/loginauth",controllers.LoginAuth)
	return r
}


