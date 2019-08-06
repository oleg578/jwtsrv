package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"../config"
	"../user"
	"github.com/google/uuid"
	"github.com/oleg578/jwts"
)

//Authorize route
// input POST JSON newuser
// input params
// email, passwd, uip (user ip)
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
	eml := r.Form.Get("email")
	if len(eml) == 0 {
		err := fmt.Errorf("wrong email")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}

	u, uerr := user.GetByEmail(eml)
	if uerr != nil {
		ResponseBuild(w, APIResp{Response: "", Error: uerr.Error()})
		return
	}
	uip := r.Form.Get("uip")

	tm := time.Now()
	texp := tm.Add(time.Minute * config.AccessDuration)
	tref := tm.Add(time.Minute * config.RefreshDuration)
	payload := make(map[string]interface{})
	payload["uid"] = u.ID
	payload["uip"] = uip
	payload["exp"] = texp.Unix()
	for _, c := range u.Claims {
		if c.AppID == r.Form.Get("appid") {
			payload["clm"] = c
		}
	}
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
