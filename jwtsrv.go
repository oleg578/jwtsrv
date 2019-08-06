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
	createHdlr := http.HandlerFunc(router.CreateHandler)
	deleteHdlr := http.HandlerFunc(router.DeleteHandler)
	searchHdlr := http.HandlerFunc(router.SearchHandler)
	listHdlr := http.HandlerFunc(router.ListHandler)

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
	//create route
	//POST source is JSON user
	mux.Handle("/create", createHdlr)
	//delete route
	//POST source is JSON user
	mux.Handle("/delete", deleteHdlr)
	//search route
	//GET param email
	mux.Handle("/search", searchHdlr)
	//list route
	//GET
	mux.Handle("/list", listHdlr)

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
