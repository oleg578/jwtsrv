package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	logger "github.com/oleg578/loglog"
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

func renderTmpl(w http.ResponseWriter, data interface{}, view string) {
	if err := TmplPool.ExecuteTemplate(w, view, data); err != nil {
		logger.Printf("view: %v template execution error: %s", view, err.Error())
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
