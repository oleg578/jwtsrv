package appregister

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/oleg578/jwtsrv/config"
)

type App struct {
	ID        string `json:"ID"`
	Resource  string `json:"Resource"`
	SecretKey string `json:"SecretKey"`
}

func GetByID(id string) (app App, err error) {
	app.ID = id
	con, errCon := redis.Dial("tcp", config.RedisDSN)
	if errCon != nil {
		return app, errCon
	}
	defer func() { _ = con.Close() }()
	rsc, errG := redis.Bytes(con.Do("HGET", "appregister", id))
	if errG != nil {
		err = fmt.Errorf("app not found: %v", errG)
		return app, err
	}
	err = json.Unmarshal(rsc, &app)
	return
}
