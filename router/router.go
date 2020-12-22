package router

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	log.Println(r.Header)
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
