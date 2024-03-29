package auth

import (
	model "async-arch/internal/domain/auth"
	"async-arch/internal/lib/httptool"
	ou "async-arch/internal/lib/osutils"
	"fmt"
	"log"
	"net/http"
)

type CheckResponse struct {
	UserUuid string `json:"uuid"`
	UserRole string `json:"role"`
}

// WithAuth - wrapper для функции обработки http-запроса, добавляющая логику авторизации по токену JWT
// В обертке выполняется проверка jwt-токена и на роли, назначенные для этого метода HTTP (роли передаются в параметре roles)
func WithAuth(handler http.HandlerFunc, roles []model.UserRole) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Получаем строку с токеном
		tokentCookie, err := r.Cookie("token")
		if err != nil {
			httptool.SetStatus401(w, "Token is not present")
			return
		}
		tokenStr := tokentCookie.Value

		if checker == nil {
			checker, err = CreateJwtTokenChecker("http", fmt.Sprintf("%s:%s", authServiceAddr, authServicePort), "GET", "/api/v1/key")
			if err != nil {
				log.Fatalln(err)
			}
		}

		checkInfo, err := checker.Check(tokenStr)
		if err != nil {
			httptool.SetStatus401(w, err.Error())
			return
		}

		// Если переданы роли, то проверяем на них
		if roles != nil {
			if len(roles) > 0 {
				ok := func(array []model.UserRole, item string) bool {
					for _, i := range array {
						if item == string(i) {
							return true
						}
					}
					return false
				}(roles, checkInfo.UserRole)
				if !ok {
					httptool.SetStatus401(w, "Unautorized user role for this method")
					return
				}
			}
		}

		// Передаем UUID пользователя и роль пользовтаеля в обработчик запроса
		r.Header.Add("X-Auth-User-UUID", checkInfo.UserUuid)
		r.Header.Add("X-Auth-User-Role", checkInfo.UserRole)

		// Вызываем исходный обработчик запроса
		handler(w, r)
	}
}

var (
	checker         *JwtTokenChecker
	authServiceAddr string = ou.GetEnvValue("AUTH_SERVER", "localhost") // Адрес сервера авторизации
	authServicePort string = ou.GetEnvValue("AUTH_SERVER_PORT", "8090") // Порт сервера авторизации
)
