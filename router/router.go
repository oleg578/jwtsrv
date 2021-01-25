package router

import (
	"encoding/json"
	"fmt"
	appreg "github.com/oleg578/jwtsrv/appregister"
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	//TmplPool templates pool
	TmplPool *template.Template
)

// APIResp response struct
type APIResp struct {
	Response interface{}
	Error    string
}

//ResponseBuild response build and send
func ResponseBuild(w http.ResponseWriter, resp APIResp) {
	if len(resp.Error) > 0 {
		time.Sleep(time.Second * 5)
	}
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	SetHeaders(w)
	_, _ = w.Write(b)
	return
}

// IndexHandler route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var Resp APIResp
	Resp.Response = r.Host
	ResponseBuild(w, Resp)
}

//SetHeaders set standard headers
func SetHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

//AppCheckMiddleware Session Middleware
func AppCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appId string
		//check header
		appId = r.Header.Get("X-AppID")
		if len(appId) == 0 {
			//get from URL query
			q := r.URL.Query()
			appId = q.Get("application_id")
		}
		if len(appId) == 0 {
			err := fmt.Errorf("wrong application id")
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		exists, errRsc := appreg.ExistsByID(appId)
		if errRsc != nil {
			ResponseBuild(w, APIResp{Response: "", Error: errRsc.Error()})
			return
		}
		if !exists {
			err := fmt.Errorf("wrong application")
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func renderTmpl(w http.ResponseWriter, data interface{}, views ...string) {
	for _, view := range views {
		if err := TmplPool.ExecuteTemplate(w, view, data); err != nil {
			log.Printf("view: %v template execution error: %s", view, err.Error())
		}
	}
}

// LoginHandler route
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//TODO add X-AppID header from get param appid
	renderTmpl(w, nil, "login.html")
}
