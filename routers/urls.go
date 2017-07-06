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
	r.HandleFunc("/redisschange", controllers.RedissDataChangeAPI)
	r.HandleFunc("/redissdel", controllers.RedissDataDelAPI)
	r.HandleFunc("/keys", controllers.Keys)
	r.HandleFunc("/keysdata", controllers.KeysDataAPI)
	r.HandleFunc("/keysdel", controllers.KeysDataDelAPI)
	r.HandleFunc("/keysexpire", controllers.KeysDataExpireAPI)
	r.HandleFunc("/keyspersist", controllers.KeysDataPersistAPI)
	//r.HandleFunc("/keysave", controllers.KeySaveAPI)
	//r.HandleFunc("/keyrename", controllers.KeyRenameAPI)
	//r.HandleFunc("/keyvaldel", controllers.KeyValDelAPI)
	//r.HandleFunc("/keydata", controllers.KeyDataAPI)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/login",controllers.Login)
	r.HandleFunc("/loginauth",controllers.LoginAuth)
	return r
}


