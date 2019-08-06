package router

import (
	"fmt"
	"log"
	"net/http"
	"time"
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
	log.Println(refreshToken)
	ResponseBuild(w, Resp)
}
