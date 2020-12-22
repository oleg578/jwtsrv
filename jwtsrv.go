package main

import (
	"log"
	"net/http"
	"time"

	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/router"
)

func main() {
	rootHandler := http.HandlerFunc(router.IndexHandler)
	authorizeHandler := http.HandlerFunc(router.AuthorizeHandler)
	renewHandler := http.HandlerFunc(router.RenewHandler)

	mux := http.NewServeMux()
	// routes
	//index route
	//GET
	mux.Handle("/", rootHandler)
	//POST
	//params apid, email, passwd
	mux.Handle("/authorize", authorizeHandler)
	//POST
	//params refresh_token
	mux.Handle("/renew", renewHandler)

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
