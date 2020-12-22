package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/oleg578/jwtsrv/user"
)

var (
	Pool *redis.Pool
)

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     300,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", 2); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}

func main() {
	Pool = newPool(":6379")
	defer func() { _ = Pool.Close() }()
	user1ID := uuid.New().String()
	user1Email := "oleg.nagornij@gmail.com"
	user1Pswd := "corner578"
	user1 := user.New(user1ID, user1Email, user1Pswd)
	asserts := make(user.AssertsMap)
	asserts["role"] = "admin"
	asserts["account"] = "*"
	claim := user.NewClaim(uuid.New().String(), asserts)
	user1.Claims = append(user1.Claims, *claim)
	if err := user1.Save(); err != nil {
		log.Fatalln("save error: ", err)
	}
	log.Printf("%+v", user1)
	//os.Exit(0)
	tUser, err := user.GetByID(user1ID)
	if err != nil {
		log.Fatalln("get user error: ", err)
	}
	log.Printf("%+v\n", tUser)
}
