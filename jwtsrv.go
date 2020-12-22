package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/oleg578/jwtsrv/config"
	"github.com/oleg578/jwtsrv/router"
	"golang.org/x/crypto/acme/autocert"
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
	//certManager
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domain),
		Cache:      autocert.DirCache(config.CertPath),
		Email:      config.AdminMail,
	}

	//server
	srv := &http.Server{
		Addr: ":https", // production
		//Addr:           ":8000", // dev
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	//https only
	log.Fatal(srv.ListenAndServeTLS("", ""))
	//local debug
	//log.Fatal(srv.ListenAndServe())
}
