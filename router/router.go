package router

import (
	"context"
	"encoding/json"
	"fmt"
	appreg "github.com/oleg578/jwtsrv/appregister"
	"github.com/oleg578/jwtsrv/logger"
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
		var (
			appId  string
			moveTo string
		)
		q := r.URL.Query()
		//check header X-AppID
		appId = r.Header.Get("X-AppID")
		if len(appId) == 0 {
			//get from URL query
			appId = q.Get("application_id")
		}
		//TODO: try get application_id from cookie
		appCookie, err := r.Cookie("app_id")
		if err != nil {
			logger.Print(err)
		} else {
			appId = appCookie.Value
		}
		if len(appId) == 0 {
			err := fmt.Errorf("wrong application id")
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		//check header X-MoveTo
		moveTo = r.Header.Get("X-MoveTo")
		if len(moveTo) == 0 {
			//get from URL query
			moveTo = q.Get("redirect_to")
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
		//propagate application_id over context
		r = r.WithContext(context.WithValue(r.Context(), "application_id", appId))
		//propagate redirect_to over context
		r = r.WithContext(context.WithValue(r.Context(), "redirect_to", moveTo))
		next.ServeHTTP(w, r)
	})
}

func renderTmpl(w http.ResponseWriter, data interface{}, view string) {
	if err := TmplPool.ExecuteTemplate(w, view, data); err != nil {
		log.Printf("view: %v template execution error: %s", view, err.Error())
	}
}

// LoginHandler route
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		AppID      string
		RedirectTo string
	}{
		r.Context().Value("application_id").(string),
		r.Context().Value("redirect_to").(string),
	}
	//set cookie with appId
	http.SetCookie(w, &http.Cookie{
		Name:  "app_id",
		Value: data.AppID,
	})
	renderTmpl(w, data, "login.html")
}
