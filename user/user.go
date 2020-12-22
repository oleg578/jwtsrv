package user

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/oleg578/jwtsrv/config"
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
	defer func() { _ = c.Close() }()
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
	con, errCon := redis.Dial("tcp", config.RedisDSN)
	if errCon != nil {
		return u, errCon
	}
	defer func() { _ = con.Close() }()
	repl, errG := redis.Bytes(con.Do("GET", id))
	if errG != nil {
		err = fmt.Errorf("user not found: %s", errG.Error())
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
	c, errCon := redis.Dial("tcp", config.RedisDSN)
	if errCon != nil {
		return u, errCon
	}
	defer func() { _ = c.Close() }()
	//get user id
	eml, errHg := c.Do("HGET", "uidbyemail", email)
	if errHg != nil {
		err = fmt.Errorf("user not found: %s", errHg.Error())
		return u, err
	}
	um, errE := redis.Bytes(c.Do("GET", eml))
	if errE != nil {
		err = fmt.Errorf("user not found: %s", errE.Error())
		return u, err
	}
	errUnm := json.Unmarshal(um, &u)
	if errUnm != nil {
		err = fmt.Errorf("user not found: %s", errUnm.Error())
	}
	return
}
