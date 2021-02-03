package router

import (
	"fmt"
	"github.com/oleg578/jwtsrv/logger"
	"github.com/oleg578/jwtsrv/token"
	"github.com/oleg578/jwtsrv/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oleg578/jwts"
	appreg "github.com/oleg578/jwtsrv/appregister"
	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/user"
)

//Authorize route
// input GET
// input params
// email, passwd
// return {"access_token":"abcd","refresh_token":"abcd"}
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		err := fmt.Errorf("wrong method")
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
	eml, pswd, errFh := formHandle(r)
	if errFh != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errFh.Error()})
		return
	}
	//redirect_to from form parse
	redirectTo := r.Form.Get("redirect_to")
	//application ID - from context or from form
	appId := r.Context().Value("application_id").(string)
	// if is a rest request, redirect may be empty
	// we check redirect only from login call
	if len(redirectTo) == 0 && len(appId) == 0 {
		err := fmt.Errorf("wrong redirect")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	if len(appId) == 0 {
		appId = r.Form.Get("appid")
	}
	if len(appId) == 0 {
		err := fmt.Errorf("wrong application")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
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
	payload["exp"] = setExpiration()

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
	//if call method is get return redirect to with jwtcode
	//if method is post return with ResponseBuild
	if r.Method == http.MethodGet {
		//generate code and save Bag
		b := &token.Bag{
			AccessToken:  AccessToken.RawStr,
			RefreshToken: RefreshToken.RawStr,
		}
		code := utils.MD5Hash(AccessToken.RawStr)
		if err := b.Save(code); err != nil {
			logger.Printf("tokens store error: %v", err)
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		//redirect
		moveTo, err := urlAddRedirect(redirectTo, code)
		if err != nil {
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		http.Redirect(w, r, moveTo, 302)
		return
	}

	if r.Method == http.MethodPost {
		ResponseBuild(w, Resp)
	}
}

func urlAddRedirect(inpath, code string) (outpath string, err error) {
	u, errParse := url.Parse(inpath)
	if errParse != nil {
		err = errParse
		return
	}
	q := u.Query()
	q.Add("code", code)
	u.RawQuery, err = url.QueryUnescape(q.Encode())
	outpath = u.String()
	return
}

func setExpiration() int64 {
	tm := time.Now()
	tref := tm.Add(time.Second * config.RefreshDuration)
	return tref.Unix()
}

func formHandle(r *http.Request) (email, passwd string, err error) {
	//get and test email
	email = r.Form.Get("email")
	if len(email) == 0 {
		err = fmt.Errorf("wrong email")
		return
	}
	//get and test password
	passwd = r.Form.Get("passwd")
	if len(passwd) == 0 {
		err = fmt.Errorf("wrong password")
		return
	}
	return
}

func payloadBuild(app appreg.App, eml, pswd string) (payload map[string]interface{}, err error) {
	payload = make(map[string]interface{})
	//try get user
	logger.Printf("user: %s, passwd: %s", eml, pswd)
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
	texp := tm.Add(time.Second * config.AccessDuration)
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
