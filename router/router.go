package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"../config"
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/google/uuid"
	"github.com/oleg578/jwts"
)

// APIResp response struct
type APIResp struct {
	Response interface{}
	Error    string
}

//ResponseBuild response build and send
func ResponseBuild(w http.ResponseWriter, resp APIResp) {
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	SetHeaders(w)
	w.Write(b)
	return
}

// IndexHandler route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var Resp APIResp
	log.Println(r.Header)
	log.Println(r.RemoteAddr) //!!!
	Resp.Response = "jwtsrv.com"
	ResponseBuild(w, Resp)
}

//Authorize route
// input POST JSON newuser
// return {"access_token":"abcd","refresh_token":"abcd"}
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != "POST" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	tm := time.Now()
	texp := tm.Add(time.Minute * config.AccessDuration)
	tref := tm.Add(time.Minute * config.RefreshDuration)
	payload := make(map[string]interface{})
	payload["uid"] = uuid.New().String()
	payload["uip"] = r.RemoteAddr
	payload["acc"] = 12018948
	payload["exp"] = texp.Unix()
	sr := strings.NewReader(payload["uid"].(string) +
		payload["uip"].(string) +
		config.SecretKey)
	jti, errjti := uuid.NewRandomFromReader(sr)
	if errjti != nil {
		jti, errjti = uuid.NewRandom()
		if errjti != nil {
			jti = uuid.New()
		}
	}
	payload["jti"] = jti.String()
	AccessToken, err := jwts.CreateTokenHS256(payload, config.SecretKey)
	if err != nil {
		log.Println(err)
	}
	payload["exp"] = tref.Unix()
	RefreshToken, errRef := jwts.CreateTokenHS256(payload, config.SecretKey)
	if errRef != nil {
		log.Println(errRef)
	}
	Resp.Response = struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken.RawStr,
		RefreshToken.RawStr,
	}
	ResponseBuild(w, Resp)
}

// CreateHandler route
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	const MAXBODYLENGTH = 1024
	var (
		Resp APIResp
	)
	if r.Method != "POST" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	Resp.Response = "OK"
	Resp.Error = ""
	ResponseBuild(w, Resp)
}

// DeleteHandler route
// input POST JSON newuser
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != "POST" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	Resp.Response = "Delete"
	Resp.Error = ""
	ResponseBuild(w, Resp)
}

// SearchHandler route
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != "GET" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	Resp.Response = "SearchHandler"
	ResponseBuild(w, Resp)
}

// ListHandler route
func ListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != "GET" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "[]", Error: err.Error()})
		return
	}
	Resp.Response = "[]"
	ResponseBuild(w, Resp)
}

//SetHeaders set standard headers
func SetHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
}

//NotFound route
func NotFound(w http.ResponseWriter, r *http.Request) {
	var Resp APIResp
	Resp.Error = "Route not found"
	ResponseBuild(w, Resp)
}
