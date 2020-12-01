package main

import (
	"flag"
	"net/http"
)

var port, secret string

func main() {
	flag.StringVar(&port, "port", "80", "http server port")
	flag.StringVar(&secret, "secret", "12345", "jwt secret")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Index)

	mux.HandleFunc("/login_jwt", LoginJwt)
	mux.HandleFunc("/profile_jwt", ProfileJwt)

	mux.HandleFunc("/login_session", LoginSessions)
	mux.HandleFunc("/profile_session", ProfileSessions)

	panic(http.ListenAndServe(":"+port, mux))
}
