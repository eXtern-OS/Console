package utils

import (
	"crypto/sha1"
	"encoding/base64"
)

func Makehash(data string) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
