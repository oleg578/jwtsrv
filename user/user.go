package user

import (
	"encoding/json"

	"../config"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type Claim struct {
	Resource string            `json:"resource"`
	Assert   map[string]string `json:"assert"`
}

type User struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Claims   []Claim `json:"claims"`
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
	return nil
}

func redisConn() (c redis.Conn, err error) {
	c, err = redis.Dial("tcp", config.RedisDSN)
	if err != nil {
		return
	}
	_, err = c.Do("SELECT", config.RedisDB)
	return c, err
}
