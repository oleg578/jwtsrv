package appregister

import (
	"fmt"

	"../config"
	"github.com/gomodule/redigo/redis"
)

type App struct {
	ID        string `json:"ID"`
	Resource  string `json:"Resource"`
	SecretKey string `json:"SecretKey"`
}

func GetByID(id string) (app App, err error) {
	app.ID = id
	c, errc := redis.Dial("tcp", config.RedisDSN)
	if errc != nil {
		return app, errc
	}
	defer c.Close()
	rsc, errG := redis.String(c.Do("HGET", "appregister", id))
	if errG != nil {
		err = fmt.Errorf("app not found")
		return app, err
	}
	app.Resource = rsc
	return
}
