package main

import "net/http"

func respondSimple(w http.ResponseWriter, code int, msg, msgType string) {
	w.Header().Add("Content-Type", "text/"+msgType+"; charset=utf-8")
	w.WriteHeader(code)
	if msg != "" {
		w.Write([]byte(msg))
	}
}
