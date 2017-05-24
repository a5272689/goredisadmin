package routers

import (
	"github.com/gorilla/mux"
	"goredisadmin/controllers"
)


func Urls() *mux.Router  {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.MainHandler)
	r.HandleFunc("/sentinels", controllers.Sentinels)
	r.HandleFunc("/sentinelsdata", controllers.SentinelsDataAPI)
	r.HandleFunc("/sentinelschange", controllers.SentinelsDataChangeAPI)
	r.HandleFunc("/sentinelsdel", controllers.SentinelsDataDelAPI)
	r.HandleFunc("/rediss", controllers.Rediss)
	r.HandleFunc("/redissdata", controllers.RedissDataAPI)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/login",controllers.Login)
	r.HandleFunc("/loginauth",controllers.LoginAuth)
	return r
}


