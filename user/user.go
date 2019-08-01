package user

import (
	"encoding/json"

	"../config"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Claim struct {
	Resource string            `json:"Resource"`
	Assert   map[string]string `json:"Assert"`
}

type User struct {
	ID       string  `json:"ID"`
	Email    string  `json:"Email"`
	Password string  `json:"Password"`
	Claims   []Claim `json:"Claims"`
}

func NewUser() *User {
	return &User{
		ID: uuid.New().String(),
	}
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
