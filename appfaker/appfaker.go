package main

import (
	"encoding/json"
	"log"

	appreg "../appregister"
	"github.com/gomodule/redigo/redis"
)

func main() {
	c, errc := redis.Dial("tcp", "192.168.1.20:6379")
	if errc != nil {
		log.Fatalln(errc)
	}
	defer c.Close()
	app := appreg.App{
		ID:        "a379ed35-a8e0-48c1-bfce-dc5eed92239c",
		Resource:  "localhost",
		SecretKey: "3dp9gudw0l19yr9ois8iu9b3220qemn8",
	}
	appS, errM := json.Marshal(app)
	if errM != nil {
		log.Fatalln(errM)
	}
	if _, err := c.Do("HSET", "appregister", app.ID, appS); err != nil {
		log.Fatalln(err)
	}

}
