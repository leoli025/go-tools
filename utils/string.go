package utils

import (
	"bytes"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandAlphaNumber(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandAlphaNumberSrc(length, str)
}

func RandAlphaNumberSrc(length int, src string) string {
	if length <= 0 {
		return ""
	}
	var buf bytes.Buffer
	var index int
	var srcLen = len(src)
	for {
		if length == 0 {
			break
		}
		index = rand.Intn(srcLen)
		buf.WriteString(src[index : index+1])
		length--
	}
	return buf.String()
}
