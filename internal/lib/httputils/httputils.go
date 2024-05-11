package httputils

import "net/http"

func SetStatus500(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}

func SetStatus401(w http.ResponseWriter, errorText string) {
	w.WriteHeader(401)
	w.Write([]byte(errorText))
}
