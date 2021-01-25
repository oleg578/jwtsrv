package router

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oleg578/jwts"
	appreg "github.com/oleg578/jwtsrv/appregister"
	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/user"
)

//Authorize route
// input POST
// input params
// email, passwd
// return {"access_token":"abcd","refresh_token":"abcd"}
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	//test method
	if r.Method != "POST" {
		err := fmt.Errorf("wrong request type")
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
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test email
	eml := r.Form.Get("email")
	if len(eml) == 0 {
		err := fmt.Errorf("wrong email")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//get and test password
	pswd := r.Form.Get("passwd")
	if len(pswd) == 0 {
		err := fmt.Errorf("wrong password")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	appId := r.Context().Value("application_id").(string)
	app, errRsc := appreg.GetByID(appId)
	if errRsc != nil {
		err := fmt.Errorf("wrong application resource")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//payload build
	payload, errPb := payloadBuild(app, eml, pswd)
	if errPb != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errPb.Error()})
		return
	}
	AccessToken, err := jwts.CreateTokenHS256(payload, app.SecretKey)
	if err != nil {
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//change expiration for refresh_token
	tm := time.Now()
	tref := tm.Add(time.Minute * config.RefreshDuration)
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

func payloadBuild(app appreg.App, eml, pswd string) (payload map[string]interface{}, err error) {
	payload = make(map[string]interface{})
	//try get user
	u, errUser := user.GetByEmail(eml)
	if errUser != nil {
		return payload, errUser
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
	payload["exp"] = texp.Unix()
	for _, c := range u.Claims {
		if c.AppID == app.ID {
			payload["clm"] = c
		}
	}
	sr := strings.NewReader(payload["uid"].(string) +
		app.SecretKey)
	jti, errJti := uuid.NewRandomFromReader(sr)
	if errJti != nil {
		jti, errJti = uuid.NewRandom()
		if errJti != nil {
			jti = uuid.New()
		}
	}
	payload["jti"] = jti.String()
	return
}
