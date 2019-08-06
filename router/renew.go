package router

import (
	"fmt"
	"net/http"
	"time"

	"../config"
	"github.com/oleg578/jwts"
	jwt "github.com/oleg578/jwts"
)

// Renew route
// input POST
// input params
// refresh_token
// return {"access_token":"abcd","refresh_token":"abcd"}
func RenewHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != "POST" {
		err := fmt.Errorf("wrong request type")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	if err := r.ParseForm(); err != nil {
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test email
	refreshToken := r.Form.Get("refresh_token")
	rtok, err := jwt.Parse(refreshToken)
	if err != nil {
		err := fmt.Errorf("token was not parsed")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	errV := rtok.Validate(config.SecretKey)
	if errV != nil {
		err := fmt.Errorf("token is not valid")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//token valid - we can build new
	payload := rtok.Payload
	tm := time.Now()
	texp := tm.Add(time.Minute * config.AccessDuration)
	tref := tm.Add(time.Minute * config.RefreshDuration)

	payload["exp"] = texp.Unix()
	AccessToken, err := jwts.CreateTokenHS256(payload, config.SecretKey)
	if err != nil {
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	payload["exp"] = tref.Unix()
	RefreshToken, errRef := jwts.CreateTokenHS256(payload, config.SecretKey)
	if errRef != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errRef.Error()})
		return
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
