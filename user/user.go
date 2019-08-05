package user

import (
	"encoding/json"

	"../config"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type AssertsMap map[string]string

type Claim struct {
	AppID    string     `json:"AppID"`
	Resource string     `json:"Resource"`
	Asserts  AssertsMap `json:"Assert"`
}

type User struct {
	ID       string  `json:"ID"`
	Email    string  `json:"Email"`
	Password string  `json:"Password"`
	Claims   []Claim `json:"Claims"`
}

func New() *User {
	return &User{
		ID: uuid.New().String(),
	}
}

func NewClaim(appid, resource string, asserts AssertsMap) *Claim {
	claim := &Claim{
		AppID:    appid,
		Resource: resource,
		Asserts:  make(AssertsMap, len(asserts)),
	}
	for key, val := range asserts {
		claim.Asserts[key] = val
	}
	return claim
}

func (u *User) Save() error {
	c, err := redisConn()
	if err != nil {
		return err
	}
	defer c.Close()
	//marshall user
	userM, errM := json.Marshal(u)
	if errM != nil {
		return errM
	}
	c.Do("SET", u.ID, userM)
	if err != nil {
		return err
	}
	err = u.EmailIndAppend()
	return err
}

func redisConn() (c redis.Conn, err error) {
	c, err = redis.Dial("tcp", config.RedisDSN)
	if err != nil {
		return
	}
	_, err = c.Do("SELECT", config.RedisDB)
	return c, err
}

func (u *User) EmailIndAppend() error {
	c, err := redisConn()
	if err != nil {
		return err
	}
	defer c.Close()
	c.Do("HSET", "uidbyemail", u.Email, u.ID)
	if err != nil {
		return err
	}
	return nil
}
