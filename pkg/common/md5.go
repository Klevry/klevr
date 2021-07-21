package common

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5Hash(buf string) string {
	hasher := md5.New()
	hasher.Write([]byte(buf))
	return hex.EncodeToString(hasher.Sum(nil))
}
