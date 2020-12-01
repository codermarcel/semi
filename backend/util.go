package main

import (
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)


func base64(input string) string {
	return strings.TrimRight(b64.URLEncoding.EncodeToString([]byte(input)), "=")
}

func base64_decode(seg  string) []byte {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	d, err := b64.URLEncoding.DecodeString(seg)
	if err != nil {
		fmt.Println("base64_decode error: " + err.Error())
		return nil
	}

	return d


	//d, err := b64.URLEncoding.DecodeString(input)
	//if err != nil {
	//	fmt.Println("base64_decode error: " + err.Error())
	//	return nil
	//}
	//return d
}

func hmac256(data string, secret string) []byte {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(data))
	return hash.Sum(nil)
}

func RenderFile(w http.ResponseWriter, r *http.Request, filename string) {
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		io.WriteString(w, "error reading html file: " +err.Error())
		return
	}

	w.Write(html)
}