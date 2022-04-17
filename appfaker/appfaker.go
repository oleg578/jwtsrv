package main

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	appreg "github.com/oleg578/jwtsrv/appregister"
)

func main() {
	c, errC := redis.Dial("tcp", ":6379")
	if errC != nil {
		log.Fatalln(errC)
	}
	defer func() { _ = c.Close() }()
	app := appreg.App{
		ID:        "a379ed35-a8e0-48c1-bfce-dc5eed92239c",
		SecretKey: "3dp9gudw0l19yr9ois8iu9b3220qemn8",
	}
	appS, errM := json.Marshal(app)
	if errM != nil {
		log.Fatalln(errM)
	}
	if _, err := c.Do("HSET",
		"appregister", app.ID, appS); err != nil {
		log.Fatalln(err)
	}
	appt := appreg.App{}
	apr, err := redis.Bytes(c.Do("HGET", "appregister", app.ID))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(string(apr))
	if errUm := json.Unmarshal(apr, &appt); errUm != nil {
		log.Fatalln(errUm)
	}
	log.Println(appt)
	log.Println(appt.SecretKey, len(appt.SecretKey))
}
