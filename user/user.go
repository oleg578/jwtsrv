package user

import (
	"encoding/json"
	"fmt"
	"strings"

	"../config"
	"github.com/gomodule/redigo/redis"
)

type AssertsMap map[string]string

type Claim struct {
	AppID   string     `json:"AppID"`
	Asserts AssertsMap `json:"Assert"`
}

type User struct {
	ID       string  `json:"ID"`
	Email    string  `json:"Email"`
	Password string  `json:"Password"`
	Claims   []Claim `json:"Claims"`
}

func New(id, email, pswd string) *User {
	return &User{
		ID:       strings.TrimSpace(id),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(pswd),
	}
}

func NewClaim(appid string, asserts AssertsMap) *Claim {
	claim := &Claim{
		AppID:   appid,
		Asserts: make(AssertsMap, len(asserts)),
	}
	for key, val := range asserts {
		claim.Asserts[key] = val
	}
	return claim
}

func (u *User) Save() error {
	c, err := redis.Dial("tcp", config.RedisDSN)
	if err != nil {
		return err
	}
	defer c.Close()
	//marshall user
	userM, errM := json.Marshal(u)
	if errM != nil {
		return errM
	}
	_, err = c.Do("SET", u.ID, userM)
	if err != nil {
		return err
	}
	err = u.EmailIndAppend(c)
	return err
}

func GetByID(id string) (u User, err error) {
	c, errc := redis.Dial("tcp", config.RedisDSN)
	if err != nil {
		return u, errc
	}
	defer c.Close()
	repl, errG := redis.Bytes(c.Do("GET", id))
	if errG != nil {
		err = fmt.Errorf("user not found")
		return u, err
	}
	//unmarshall user
	err = json.Unmarshal(repl, &u)

	return
}

func (u *User) EmailIndAppend(c redis.Conn) error {
	_, err := c.Do("HSET", "uidbyemail", u.Email, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetByEmail(email string) (u User, err error) {
	c, errc := redis.Dial("tcp", config.RedisDSN)
	if err != nil {
		return u, errc
	}
	defer c.Close()
	//get user id
	eml, errhg := c.Do("HGET", "uidbyemail", email)
	if errhg != nil {
		err = fmt.Errorf("user not found")
		return u, err
	}
	um, erre := redis.Bytes(c.Do("GET", eml))
	if erre != nil {
		err = fmt.Errorf("user not found")
		return u, err
	}
	err = json.Unmarshal(um, &u)
	if err != nil {
		err = fmt.Errorf("user not found")
	}
	return
}
