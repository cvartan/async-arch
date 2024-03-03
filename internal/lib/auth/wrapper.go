package auth

import (
	model "async-arch/internal/domain/auth"
	"async-arch/internal/lib/base"
	"async-arch/internal/lib/httptool"
	"async-arch/internal/lib/sysenv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

		// Делаем запрос на сервис авторизации - получаем расшифровку токена (к сожалению в версии 1.22 что-то сломали в декодировании PEM-формата и ключ теперь получить нормально нельзя - он не преобразуется из строки в структуру rsa.PublicKey)
		resp, err := base.App.Request(keyRequestId, []byte(tokenStr), nil, nil, nil, nil)
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}
		defer resp.Body.Close()

		var checkInfo CheckResponse
		source, err := io.ReadAll(resp.Body)
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}
		err = json.NewDecoder(strings.NewReader(string(source))).Decode(&checkInfo)
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}

		/*
			TODO: разобраться с декодером PEM-формата и сделать через локальную проверку (или найти другие библиотеки по работе с RSA)
		*/

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

var authServiceAddr string = sysenv.GetEnvValue("AUTH_SERVER", "localhost") // Адрес сервера авторизации
var authServicePort string = sysenv.GetEnvValue("AUTH_SERVER_PORT", "8090") // Порт сервера авторизации

var keyRequestId string = uuid.NewString() // ИД шаблона для запроса проверки токена

func init() {
	// Добавляем шаблон запроса проверки токена
	base.App.AddPostRequest(keyRequestId, fmt.Sprintf("http://%s:%s", authServiceAddr, authServicePort), "/api/v1/check")
}
