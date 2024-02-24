package auth

import (
	"async-arch/lib/base"
	"async-arch/util"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

var authServiceAddr string = util.GetEnvValue("AUTH_SERVER", "localhost")
var authServicePort string = util.GetEnvValue("AUTH_SERVER_PORT", "8090")

var loginRequestId string = uuid.NewString()
var authRequestId string = uuid.NewString()

func WrapAuth(handler func(response http.ResponseWriter, request *http.Request), permission string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionId := r.Header.Get("X-Session-ID")
		if sessionId == "" {
			resp, err := base.App.Request(loginRequestId, nil, nil, nil, map[string]string{"Authorization": r.Header.Get("Authorization")})
			if err != nil {
				w.WriteHeader(401)
				w.Write([]byte(err.Error()))
				return
			}
			if resp.StatusCode == 401 {
				w.WriteHeader(401)
				return
			}
			sessionId = resp.Header.Get("X-Session-ID")
		}

		requestStr := fmt.Sprintf("{\"session\":\"%s\",\n\"permission\":\"%s\"}", sessionId, permission)

		resp, err := base.App.Request(authRequestId, []byte(requestStr), nil, nil, nil)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}
		if resp.StatusCode == 401 {
			w.WriteHeader(401)
			return
		}

		handler(w, r)
	}
}

func init() {
	base.App.AddPostRequest(loginRequestId, fmt.Sprintf("http://%s:%s", authServiceAddr, authServicePort), "/api/v1/login/")
	base.App.AddPostRequest(authRequestId, fmt.Sprintf("http://%s:%s", authServiceAddr, authServicePort), "/api/v1/auth")
}
