package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/oleg578/jwtsrv/appflags"
	"github.com/oleg578/jwtsrv/config"
	logger "github.com/oleg578/loglog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/oleg578/jwtsrv/router"
)

func main() {
	appflags.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	appflags.Dev, _ = strconv.ParseBool(os.Getenv("DEV"))
	setEnv(appflags.Dev)
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
	if !appflags.Dev {
		srv.Addr = ":https"
		srv.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
		}
	}
	if appflags.Dev {
		srv.Addr = ":http"
	}
	//production
	if !appflags.Dev {
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		logger.Fatal(srv.ListenAndServeTLS("", ""))
	}
	//local debug
	if appflags.Dev {
		log.Fatal(srv.ListenAndServe())
	}
}

func setEnv(devmode bool) {
	if devmode {
		config.AdminMail = "oleg.nagornij@gmail.com"
		config.Domain = "accounts.bwretail.com"
		config.CertPath = "/etc/autocert/ssl/"
		config.CODELIFETIME = 900
		config.RedisDSN = `127.0.0.1:6379`
		config.TemplateDir = "./tmpl/"
		config.LogPath = "./log/jwtsrv.log"
	} else {
		config.AdminMail = "oleg.nagornij@gmail.com"
		config.Domain = "accounts.bwretail.com"
		config.CertPath = "/etc/autocert/ssl/"
		config.CODELIFETIME = 900
		config.RedisDSN = `127.0.0.1:6379`
		config.TemplateDir = "/var/www/tmpl/"
		config.LogPath = "/var/log/jwtsrv.log"
	}

}
