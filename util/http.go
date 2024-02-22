package util

import "net/http"

func GenerateHTTP500Response(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
