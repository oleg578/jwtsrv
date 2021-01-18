package main

import (
	"github.com/oleg578/jwtsrv/config"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oleg578/jwtsrv/router"
)

func main() {

	rootHandler := http.HandlerFunc(router.IndexHandler)
	authorizeHandler := http.HandlerFunc(router.AuthorizeHandler)
	renewHandler := http.HandlerFunc(router.RenewHandler)
	loginHandler := http.HandlerFunc(router.LoginHandler)

	mux := http.NewServeMux()
	// routes
	//index route
	//GET
	mux.Handle("/", router.AppCheckMiddleware(rootHandler))
	//GET
	mux.Handle("/login", router.AppCheckMiddleware(loginHandler))
	//POST
	//params apid, email, passwd
	mux.Handle("/authorize", router.AppCheckMiddleware(authorizeHandler))
	//POST
	//params refresh_token
	mux.Handle("/renew", router.AppCheckMiddleware(renewHandler))

	//server
	//certManager
	//certManager := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist(config.Domain),
	//	Cache:      autocert.DirCache(config.CertPath),
	//	Email:      config.AdminMail,
	//}

	//server
	srv := &http.Server{
		//Addr: ":https", // production
		Addr:           ":5000", // dev
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		//TLSConfig: &tls.Config{
		//	GetCertificate: certManager.GetCertificate,
		//},
	}
	//https only
	//log.Fatal(srv.ListenAndServeTLS("", ""))
	//local debug
	log.Fatal(srv.ListenAndServe())
}

func init() {
	if err := os.Chdir(config.TemplateDir); err != nil {
		log.Fatal(err)
	}

	router.TmplPool = template.Must(template.ParseGlob(`*.html`))
}
