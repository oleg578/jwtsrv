package main

import (
	"log"

	"../user"
)

func main() {
	me := user.New()
	log.Printf("%+v", me)
}
