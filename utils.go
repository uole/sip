package sip

import "crypto/md5"

func MD5(b []byte) []byte {
	hash := md5.New()
	hash.Write(b)
	return hash.Sum(nil)
}
