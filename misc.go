package main

import (
	"crypto/md5"
	"fmt"
)

func MD5Encode(s string) string {
	data := []byte(s)
	buff := md5.Sum(data)
	return fmt.Sprintf("%x", buff)
}
