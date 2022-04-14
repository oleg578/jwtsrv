package main

import (
	"fmt"
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
			if _, err := c.Do("SELECT", 0); err != nil {
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
	nick := "Oleh"
	user1Pswd := "corner578"
	secret := "secret"
	user1 := user.New(user1ID, user1Email, nick, user1Pswd, secret)
	claim := user.NewClaim("*", "admin")
	user1.Claims = append(user1.Claims, *claim)
	claim1 := user.NewClaim("localhost", "superadmin")
	user1.Claims = append(user1.Claims, *claim1)
	if err := user1.Save(); err != nil {
		log.Fatalln("save error: ", err)
	}
	fmt.Printf("%+v\n", user1)
}
