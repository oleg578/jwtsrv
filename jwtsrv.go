package main

import (
	"log"
	"net/http"
	"time"

	"./config"
	"./router"
)

func main() {
	rootHdlr := http.HandlerFunc(router.IndexHandler)
	authorizeHdlr := http.HandlerFunc(router.AuthorizeHandler)
	renewHdlr := http.HandlerFunc(router.RenewHandler)

	mux := http.NewServeMux()
	// routes
	//index route
	//GET
	mux.Handle("/", rootHdlr)
	//POST
	//params apid, email, passwd
	mux.Handle("/authorize", authorizeHdlr)
	//POST
	//params refresh_token
	mux.Handle("/renew", renewHdlr)

	//server
	//for cluster
	srv := &http.Server{
		Addr:           config.Host + ":5000",
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
