package controllers



import (
	"net/http"
	"github.com/goincremental/negroni-sessions"
	//"fmt"
	//"time"
	//"encoding/json"
	"github.com/flosch/pongo2"
	"goredisadmin/models"
	//"goredisadmin/utils"
	//"github.com/bitly/go-simplejson"
	//"strconv"
	//"strings"
	//"io/ioutil"
	//"goredisadmin/utils"
	"fmt"
)

func initconText(r *http.Request) (pongo2.Context) {
	session := sessions.GetSession(r)
	redisClient,err:=models.RedisPool.Get()
	defer models.RedisPool.Put(redisClient)
	if err!=nil{
		return pongo2.Context{}
	}
	sentinels_keys,_:=redisClient.Hkeys("goredisadmin:sentinels:hash")
	redis_keys,_:=redisClient.Hkeys("goredisadmin:rediss:hash")
	user:=session.Get("user")
	username:=session.Get("username")
	if username==nil{
		username=user
	}
	userrole:=session.Get("role")
	if userrole==nil{
		userrole=""
	}
	return pongo2.Context{"username":username,"userrole":userrole,"urlpath":r.URL.Path,"sentinels":len(sentinels_keys),"redis":len(redis_keys)}
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	//http.SetCookie(w,&http.Cookie{Name:"csrftoken",Value:string(time.Now().String()),MaxAge:60})
	//fmt.Println(r.URL.Path)
	//fmt.Println(utils.ConfLoad())
	//fmt.Fprintln(w, session.Get("user"))
	//userdb:=&models.User{UserName:"jkljdaklsjfkl"}
	//dbpass,err:=userdb.GetPassWord()
	//fmt.Println(dbpass,err)
	//http.Redirect(w,r,"/",http.StatusFound)
	tpl,err:=pongo2.FromFile("views/contents/index.html")
	tpl = pongo2.Must(tpl,err)
	fmt.Println(session.Get("user"))
	tpl.ExecuteWriter(initconText(r), w)
}
