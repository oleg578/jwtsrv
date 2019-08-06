package router

import (
	"fmt"
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
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test email
	eml := r.Form.Get("email")
	if len(eml) == 0 {
		err := fmt.Errorf("wrong email")
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test password
	pswd := r.Form.Get("passwd")
	if len(pswd) == 0 {
		err := fmt.Errorf("wrong password")
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test user IP
	uip := r.Form.Get("uip")
	if len(uip) == 0 {
		err := fmt.Errorf("user ip detection error")
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test appid
	appid := r.Form.Get("appid")
	if len(appid) == 0 {
		err := fmt.Errorf("app ip error")
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//payload build
	payload, errpb := payloadBuild(appid, eml, pswd, uip)
	if errpb != nil {
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: errpb.Error()})
		return
	}
	AccessToken, err := jwts.CreateTokenHS256(payload, config.SecretKey)
	if err != nil {
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//change expiration for refresh_token
	tm := time.Now()
	tref := tm.Add(time.Minute * config.RefreshDuration)
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

func payloadBuild(appid, eml, pswd, uip string) (payload map[string]interface{}, err error) {
	payload = make(map[string]interface{})
	//try get user
	u, uerr := user.GetByEmail(eml)
	if uerr != nil {
		return payload, uerr
	}
	//test user passwd
	if u.Password != pswd {
		err = fmt.Errorf("wrong password")
		return
	}
	tm := time.Now()
	texp := tm.Add(time.Minute * config.AccessDuration)
	payload["uid"] = u.ID
	payload["eml"] = u.Email
	payload["uip"] = uip
	payload["exp"] = texp.Unix()
	for _, c := range u.Claims {
		if c.AppID == appid {
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
	return
}
