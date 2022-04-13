package router

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/oleg578/jwtsrv/appflags"
	logger "github.com/oleg578/loglog"

	"github.com/google/uuid"
	"github.com/oleg578/jwts"
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
	if appflags.Debug {
		log.Printf("AuthorizeHandler request: %+v\n", r)
	}
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
	u, errURL := url.Parse(r.Referer())
	if appflags.Debug {
		log.Printf("request referer: %s", r.Referer())
		log.Printf("referer parsed: %+v", u)
	}
	if errURL != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errURL.Error()})
		return
	}
	ref := u.Host
	if appflags.Debug {
		log.Printf("referer: %s", ref)
	}
	if err := r.ParseForm(); err != nil {
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	if appflags.Debug {
		log.Printf("form: %+v\n", r.Form)
	}
	eml, pswd, errFh := formHandle(r)
	if errFh != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errFh.Error()})
		return
	}
	if appflags.Debug {
		log.Printf("email: %s, passwd: %s", eml, pswd)
	}

	//payload build
	payload, secret, errPb := payloadBuild(ref, eml, pswd)
	if errPb != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errPb.Error()})
		return
	}
	AccessToken, err := jwts.CreateTokenHS256(payload, secret)
	if err != nil {
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	//change expiration for refresh_token
	payload["exp"] = setExpiration()
	RefreshToken, errRef := jwts.CreateTokenHS256(payload, secret)
	if errRef != nil {
		ResponseBuild(w, APIResp{Response: "", Error: errRef.Error()})
		return
	}
	//if call method is get return redirect to with jwtcode
	//if method is post return with ResponseBuild
	if r.Method == http.MethodPost {
		Resp.Response = struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			AccessToken.RawStr,
			RefreshToken.RawStr,
		}
		ResponseBuild(w, Resp)
	}
	if r.Method == http.MethodGet {
		//redirect_to from form parse
		redirectTo := r.Form.Get("redirect_to")
		//redirect to client with jwtcode
		redirectURL, err := urlAddParam(redirectTo, "access_token", AccessToken.RawStr)
		if err != nil {
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		redirectURL, errR := urlAddParam(redirectURL, "refresh_token", RefreshToken.RawStr)
		if errR != nil {
			ResponseBuild(w, APIResp{Response: "", Error: errR.Error()})
			return
		}
		// if is a rest request, redirect may be empty
		// we check redirect only from login call
		if len(redirectTo) == 0 {
			err := fmt.Errorf("wrong redirect")
			ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
			return
		}
		http.Redirect(w, r, redirectURL, 302)
		return
	}

}

func urlAddParam(inpath, name, value string) (outpath string, err error) {
	u, errParse := url.Parse(inpath)
	if errParse != nil {
		err = errParse
		return
	}
	q := u.Query()
	q.Add(name, value)
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

func payloadBuild(referer, eml, pswd string) (payload map[string]interface{}, secret string, err error) {
	payload = make(map[string]interface{})
	//try to get user
	logger.Printf("user: %s, passwd: %s", eml, pswd)
	u, err := user.GetByEmail(eml)
	if err != nil {
		return nil, "", err
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
		if c.Resource == referer {
			payload["clm"] = c
		}
	}
	sr := strings.NewReader(payload["uid"].(string) + u.SecretKey)
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
