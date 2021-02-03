package config

const (
	AdminMail = "oleg.nagornij@gmail.com"
	Domain    = "accounts.bwretail.com"
	CertPath  = "/etc/autocert/ssl/"

	MAXBODYLENGTH = 4096

	AccessDuration  = 15 * 60      // 15 minutes (in seconds)
	RefreshDuration = 12 * 60 * 60 // 12 hour (in seconds)

	CODELIFETIME = 900
)

var (
	RedisDSN      = `127.0.0.1:6379`
	RedisDSNLocal = `192.168.1.20:6379`

	TemplateDirLocal = "./tmpl/"
	TemplateDir      = "/var/www/tmpl/"

	LogPathLocal = "./log/jwtsrv.log"
	LogPath      = "/var/log/jwtsrv.log"
)
