package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5Hash calc MD5 HASH
func MD5Hash(s string) string {
	md5sum := md5.Sum([]byte(s))
	return hex.EncodeToString(md5sum[:])
}
