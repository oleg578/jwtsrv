package config

const (
	AdminMail = "oleg.nagornij@gmail.com"
	Domain    = "accounts.bwretail.com"
	CertPath  = "/etc/autocert/ssl/"

	MAXBODYLENGTH = 2048

	//RedisDSN = `192.168.1.20:6379`
	RedisDSN = `127.0.0.1:6379`

	AccessDuration  = 1440 * 1000  // 24 hour
	RefreshDuration = 43200 * 1000 // 30*24 hour

	//TemplateDir = "./tmpl/"
	TemplateDir = "/var/www/tmpl/"

	//LogPath = "./log/jwtsrv.log"
	LogPath = "/var/log/jwtsrv.log"
)
