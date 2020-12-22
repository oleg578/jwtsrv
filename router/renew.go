package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oleg578/jwts"
	appreg "github.com/oleg578/jwtsrv/appregister"
	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/user"
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
	//test appid in white list from header Bw-Appid
	appid := r.Header.Get("Bw-Appid")
	if len(appid) == 0 {
		err := fmt.Errorf("wrong application resource")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	app, errRsc := appreg.GetByID(appid)
	if errRsc != nil {
		err := fmt.Errorf("wrong application resource")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//prevent maliciously sending
	if r.ContentLength > config.MAXBODYLENGTH {
		err := fmt.Errorf("request body length limit exceeded")
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
	refTok, err := jwts.Parse(refreshToken)
	if err != nil {
		err := fmt.Errorf("token was not parsed")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	errV := refTok.Validate(app.SecretKey)
	if errV != nil {
		err := fmt.Errorf("token is not valid")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//token valid - we can build new
	//test user exists
	//try get user
	_, errUser := user.GetByID(refTok.Payload["uid"].(string))
	if errUser != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errUser.Error()})
		return
	}
	//build tokens
	payload := refTok.Payload
	tm := time.Now()
	texp := tm.Add(time.Minute * config.AccessDuration)
	tref := tm.Add(time.Minute * config.RefreshDuration)
	payload["exp"] = texp.Unix()
	AccessToken, err := jwts.CreateTokenHS256(payload, app.SecretKey)
	if err != nil {
		time.Sleep(time.Second * 5)
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	payload["exp"] = tref.Unix()
	RefreshToken, errRef := jwts.CreateTokenHS256(payload, app.SecretKey)
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
