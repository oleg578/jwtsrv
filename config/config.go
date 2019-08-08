package config

const (
	MAXBODYLENGTH = 2048

	//Host = "10.132.146.197"
	Host = ""

	RedisDSN = `192.168.1.20:6379`
	RedisDB  = 2

	SecretKey       = "3dp9gudw0l19yr9ois8iu9b3220qemn8"
	AccessDuration  = 1440  // 24 hour
	RefreshDuration = 43200 // 30*24 hour
)
