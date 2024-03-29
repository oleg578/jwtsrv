package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/oleg578/jwtsrv/config"
	logger "github.com/oleg578/loglog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/oleg578/jwtsrv/router"
)

func main() {

	//if local development mode DEVMODE=true
	DevMode, _ := strconv.ParseBool(os.Getenv("DEVMODE"))
	Production := !DevMode
	if !Production {
		config.RedisDSN = config.RedisDSNLocal
		config.TemplateDir = config.TemplateDirLocal
		config.LogPath = config.LogPathLocal
	}
	//logger
	if err := logger.New(config.LogPath, "", logger.LstdFlags); err != nil {
		log.Fatal(err)
	}
	//templates
	router.TmplPool = template.Must(template.ParseGlob(config.TemplateDir + "*.html"))

	//handlers
	rootHandler := http.HandlerFunc(router.IndexHandler)
	authorizeHandler := http.HandlerFunc(router.AuthorizeHandler)
	renewHandler := http.HandlerFunc(router.RenewHandler)
	loginHandler := http.HandlerFunc(router.LoginHandler)
	originHandler := http.HandlerFunc(router.OriginHandler)

	mux := http.NewServeMux()
	// routes
	//index route
	//GET
	//mux.Handle("/", router.AppCheckMiddleware(rootHandler))
	mux.Handle("/", rootHandler)
	//GET
	mux.Handle("/login", router.AppCheckMiddleware(loginHandler))
	//GET
	//params apid, email, passwd
	mux.Handle("/authorize", router.AppCheckMiddleware(authorizeHandler))
	//POST
	//params refresh_token
	mux.Handle("/renew", router.AppCheckMiddleware(renewHandler))

	//GET tokens pair for code
	mux.Handle("/origin", router.AppCheckMiddleware(originHandler))

	//server certManager
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(config.Domain),
		Cache:      autocert.DirCache(config.CertPath),
		Email:      config.AdminMail,
	}
	//server
	srv := &http.Server{
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if Production {
		srv.Addr = ":https"
		srv.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
	}
	if DevMode {
		srv.Addr = ":5000"
	}
	//production
	if Production {
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		logger.Fatal(srv.ListenAndServeTLS("", ""))
	}
	//local debug
	if DevMode {
		log.Fatal(srv.ListenAndServe())
	}
}
