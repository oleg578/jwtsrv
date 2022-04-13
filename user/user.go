package user

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/oleg578/jwtsrv/config"
)

type AssertsMap map[string]string

// User can have asserts for each Application ID
type Claim struct {
	Resource string     `json:"Resource"`
	Asserts  AssertsMap `json:"Assert"` // map[string]string resource -> role
}

type User struct {
	ID        string  `json:"ID"`
	Email     string  `json:"Email"`
	Nickname  string  `json:"Nickname"`
	Password  string  `json:"Password"`
	SecretKey string  `json:"SecretKey"`
	Claims    []Claim `json:"Claims"` //claims for each application
}

func New(id, email, nickname, pswd, secret string) *User {
	return &User{
		ID:        strings.TrimSpace(id),
		Email:     strings.TrimSpace(email),
		Nickname:  strings.TrimSpace(nickname),
		Password:  strings.TrimSpace(pswd),
		SecretKey: strings.TrimSpace(secret),
	}
}

func NewClaim(resource string, asserts AssertsMap) *Claim {
	claim := &Claim{
		Resource: resource,
		Asserts:  make(AssertsMap, len(asserts)),
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
	err = u.EmailIndexAppend(c)
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

func (u *User) EmailIndexAppend(c redis.Conn) error {
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
