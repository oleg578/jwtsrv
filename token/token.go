package token

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/oleg578/jwtsrv/config"
)

//expiration 15 minutes
type Bag struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (b *Bag) Save(code string) error {
	var err error
	c, errDial := redis.Dial("tcp", config.RedisDSN)
	if errDial != nil {
		return errDial
	}
	defer func() { _ = c.Close() }()
	//marshall bag
	bM, errM := json.Marshal(b)
	if errM != nil {
		return errM
	}
	_, err = c.Do("SET", code, bM, config.CODELIFETIME)
	return err
}

func Get(code string) (b Bag, err error) {
	con, errCon := redis.Dial("tcp", config.RedisDSN)
	if errCon != nil {
		return b, errCon
	}
	defer func() { _ = con.Close() }()
	repl, errG := redis.Bytes(con.Do("GET", code))
	if errG != nil {
		err = fmt.Errorf("bag not found or expired: %s", errG.Error())
		return b, err
	}
	//unmarshall bag
	err = json.Unmarshal(repl, &b)
	return
}
