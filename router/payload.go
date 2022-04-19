package router

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oleg578/jwtsrv/appflags"
	"github.com/oleg578/jwtsrv/appregister"
	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/user"
	logger "github.com/oleg578/loglog"
)

func payloadBuild(referer, eml, pswd, appid string) (payload map[string]interface{}, secret string, err error) {
	payload = make(map[string]interface{})
	//try to get user
	if appflags.Debug {
		logger.Printf("user: %s, passwd: %s", eml, pswd)
	}
	u, err := user.GetByEmail(eml)
	if err != nil {
		return nil, "", err
	}
	//test user passwd
	if u.Password != pswd {
		err = fmt.Errorf("wrong password")
		return
	}
	//check appid
	app, err := appregister.GetByID(appid)
	if err != nil {
		return nil, "", err
	}
	secret = app.SecretKey

	tm := time.Now()
	texp := tm.Add(time.Second * config.AccessDuration)
	payload["uid"] = u.ID
	payload["eml"] = u.Email
	payload["nick"] = u.Nickname
	payload["exp"] = texp.Unix()
	payload["role"] = "guest"

	for _, c := range u.Claims {
		if appflags.Debug {
			logger.Printf("claim: %+v", c)
		}
		if c.Resource == referer {
			payload["role"] = c.Role
			break
		}
		if c.Resource == "*" {
			payload["role"] = c.Role
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
