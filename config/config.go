package config

import "time"

var (
	AdminMail = "oleg.nagornij@gmail.com"
	Domain    = "accounts.bwretail.com"
	CertPath  = "/etc/autocert/ssl/"

	MAXBODYLENGTH int64 = 4096

	AccessDuration  time.Duration = 15 * 60      // 15 minutes (in seconds)
	RefreshDuration time.Duration = 12 * 60 * 60 // 12 hour (in seconds)

	CODELIFETIME = 900

	RedisDSN = `127.0.0.1:6379`

	TemplateDirLocal = "./tmpl/"
	TemplateDir      = "/var/www/tmpl/"

	LogPathLocal = "./log/jwtsrv.log"
	LogPath      = "/var/log/jwtsrv.log"
)
