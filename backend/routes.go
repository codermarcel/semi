package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/sessions"
)


func Index(w http.ResponseWriter, r *http.Request) {
	RenderFile(w, r, "./html/index.html")
	return
}


func LoginJwt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")

	if username == "" && password == "" {
		RenderFile(w, r, "./html/login_jwt.html")
		return
	}

	if username != "username" {
		io.WriteString(w, "bad username")
		return
	}

	if password != "password" {
		io.WriteString(w, "bad password")
		return
	}

	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	payload := make(Payload)
	payload["sub"] = "1" //user id
	payload["exp"] = time.Now().Add(time.Hour * 2).Unix() //expiration
	payload["email"] = "user@email.com" //user email
	payload["username"] = "user" //user email
	payload["role"] = "user" //user email

	jsonHeader := header.ToJson()
	jsonPayload := payload.ToJson()

	base64Header := base64(jsonHeader)
	base64Payload := base64(jsonPayload)

	signature := hmac256(base64Header + "." + base64Payload, secret)

	io.WriteString(w, fmt.Sprintf("%s.%s.%s", base64Header, base64Payload, base64(string(signature))))
}


func ProfileJwt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	jwt := query.Get("jwt")

	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		io.WriteString(w, "bad jwt format")
		return
	}

	jsonHeader := base64_decode(parts[0])
	jsonPayload := base64_decode(parts[1])
	signature := base64_decode(parts[2])

	realSignature := hmac256(parts[0] + "." + parts[1], secret)

	if bytes.Compare(signature, realSignature) != 0{
		io.WriteString(w, "signature does not match")
		return
		io.WriteString(w, fmt.Sprintf("%s\n%s\n%s\n\n%s\n%s\n", parts[0], parts[1], parts[2], realSignature, jsonHeader))
	}

	payload := make(Payload)
	err := payload.LoadJSON(jsonPayload)
	if err != nil {
		io.WriteString(w, "invalid json payload, err:" + err.Error())
		return
	}


	data :=  &ProfileData{payload.Get("username"), payload.Get("role"), "jwt"}

	tmpl := template.Must(template.ParseFiles("./html/profile.html"))
	tmpl.Execute(w,data)
}



var store = sessions.NewCookieStore([]byte(secret))


func LoginSessions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	username := query.Get("username")
	password := query.Get("password")

	if username == "" && password == "" {
		RenderFile(w, r, "./html/login_session.html")
		return
	}

	//select * from users where username = ?
	if username != "username" {
		io.WriteString(w, "bad username")
		return
	}

	if password != "password" {
		io.WriteString(w, "bad password")
		return
	}

	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	role := "user"
	session.Values["username"] = username
	session.Values["role"] = role

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("./html/profile.html"))
	data := &ProfileData{username, role, "session"}
	tmpl.Execute(w, data)
}

type ProfileData struct {
	Username string
	Role string
	Method string
}

func ProfileSessions(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := session.Values["username"]
	role := session.Values["role"]

	if username == nil {
		io.WriteString(w, "not authorized")
		return
	}


	tmpl := template.Must(template.ParseFiles("./html/profile.html"))
	tmpl.Execute(w, &ProfileData{username.(string), role.(string), "session"})
}

