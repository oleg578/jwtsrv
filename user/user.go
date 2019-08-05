package user

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
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

func New(id, email, pswd string) *User {
	return &User{
		ID:       strings.TrimSpace(id),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(pswd),
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

func (u *User) Save(pool *redis.Pool) error {
	c := pool.Get()
	defer c.Close()
	//marshall user
	userM, errM := json.Marshal(u)
	if errM != nil {
		return errM
	}
	_, err := c.Do("SET", u.ID, userM)
	if err != nil {
		return err
	}
	err = u.EmailIndAppend(pool)
	return err
}

func GetByID(id string, pool *redis.Pool) (u User, err error) {
	c := pool.Get()
	defer c.Close()
	repl, errG := redis.Bytes(c.Do("GET", id))
	if errG != nil {
		return u, errG
	}
	log.Println(string(repl))
	//unmarshall user
	errUM := json.Unmarshal(repl, &u)
	return u, errUM
}

func (u *User) EmailIndAppend(pool *redis.Pool) error {
	c := pool.Get()
	defer c.Close()
	_, err := c.Do("HSET", "uidbyemail", u.Email, u.ID)
	if err != nil {
		return err
	}
	return nil
}
