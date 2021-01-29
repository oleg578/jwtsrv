package router

import (
	"fmt"
	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/token"
	"net/http"
)

// Origin route
// return bag of tokens by code
// Method GET
func OriginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		Resp APIResp
	)
	if r.Method != http.MethodGet {
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
	q := r.URL.Query()
	code := q.Get("code")
	if len(code) == 0 {
		err := fmt.Errorf("wrong parameters")
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	b, err := token.Get(code)
	if err != nil {
		ResponseBuild(w, APIResp{Response: "", Error: err.Error()})
		return
	}
	Resp.Response = struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		b.AccessToken,
		b.RefreshToken,
	}
	ResponseBuild(w, Resp)
}
