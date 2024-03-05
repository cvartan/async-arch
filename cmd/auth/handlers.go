// Обработчики http-запросов сервиса авторизации

package main

import (
	model "async-arch/internal/domain/auth"
	authlib "async-arch/internal/lib/auth"
	"async-arch/internal/lib/base"
	"async-arch/internal/lib/event"
	"async-arch/internal/lib/httptool"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var authPrivateKey *rsa.PrivateKey //Приватный ключ для подписи токена JWT

// Инициализация API сервиса авторизации
func initHandlers() {
	base.App.HandleFunc("POST /api/v1/users", handleCreateUser)
	base.App.HandleFunc("POST /api/v1/login", handleAuthentificate)
	base.App.HandleFunc("GET /api/v1/key", handleGetKey)
	base.App.HandleFunc("POST /api/v1/check", handleCheck)

	// Для простоты генерируем приватный ключ (вместо использования заранее сгененренных ключей)
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
	}
	authPrivateKey = pk
}

// Обработка запроса добавления пользователя
func handleCreateUser(w http.ResponseWriter, r *http.Request) {

	// Получаем тело запроса добавления пользователя
	var userRq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRq); err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	// Формируем сущность "Пользователь"
	user := &model.User{
		Beak:     model.ParrotBeak(userRq.Beak),
		Name:     userRq.Name,
		EMail:    userRq.EMail,
		Password: userRq.Password,
		Uuid:     uuid.NewString(),
		Role:     model.UserRole(userRq.Role),
	}

	// Сохраняем данные пользователя в БД
	repo, _ := base.App.GetDomainRepository("auth")

	err := repo.Append(user)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	// Формируем данные для ответа
	userResp := &UserResponse{
		ID:    user.ID,
		Uuid:  user.Uuid,
		Beak:  string(user.Beak),
		Name:  user.Name,
		Email: user.EMail,
		Role:  string(user.Role),
	}

	if err := json.NewEncoder(w).Encode(userResp); err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	// Возвращаем статус 201
	w.WriteHeader(201)

	// Отправляем событие в очередь
	eventData := UserEventData{
		Uuid:  user.Uuid,
		Name:  user.Name,
		Email: user.EMail,
		Role:  string(user.Role),
	}
	evnt, err := eventProducer.ProduceEventData(event.AUTH_CUD_USER_CREATED, user.Uuid, reflect.TypeOf(*user).String(), eventData)
	if err != nil {
		log.Fatal(err)
	}
	// Сохраняем событие в БД
	e := &Event{
		Event:         *evnt,
		UserEventData: eventData,
	}

	err = repo.Append(e)
	if err != nil {
		log.Fatal(err)
	}
}

// Метод аутентификации попугая
func handleAuthentificate(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Получаем уникальный оттиск ключа в качестве логина и пароль
	beak, password, ok := r.BasicAuth()
	if !ok {
		httptool.SetStatus401(w, "Expected basic authorization")
		return
	}
	if beak == "" {
		httptool.SetStatus401(w, "Empty parrot beak profile")
		return
	}
	if password == "" {
		httptool.SetStatus401(w, "Empty password")
		return
	}

	// Ищем попугая с таким оттистком клюва
	var user model.User
	repo, _ := base.App.GetDomainRepository("auth")
	err := repo.Get(&user, map[string]interface{}{"beak": beak})
	if err != nil {
		httptool.SetStatus401(w, "Unregistered user")
		return
	}

	// Определяем хэш MD5 для пароля (так как в БД вместо пароля хранится хэш)
	hash := md5.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil))

	// Проверяем пароль
	if password != user.Password {
		httptool.SetStatus401(w, "Incorrect password")
		return
	}

	// Формирование JWT токена
	// Устанавливаем срок действия - 5 минут с момента создания
	expiresAt := time.Now().Add(time.Minute * 5)

	// Задаем структуру токена JWT
	claims := authlib.AuthClaims{
		UserUuid: user.Uuid,
		UserRole: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth",
			Subject:   user.Name,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Собираем токен
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Подписываем токен (используем алгоритм RSA)
	tokenStr, err := token.SignedString(authPrivateKey)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	// Возвращаем токен в заголовке (потому что у меня кривые инструменты тестирования API)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenStr,
		Expires: expiresAt,
	})

	// Записываем статус 200
	w.WriteHeader(200)
}

// Метод получения публичного ключа сервиса авторизации
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text")
	// Формируем текст в формате PEM
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&authPrivateKey.PublicKey),
	}
	b := pem.EncodeToMemory(block)
	w.Write(b)
}

// Метод проверки токена (добавляем так как, судя по всему, в 1.22 сломали декодер PEM-формата - поэтому публичный ключ передать не получается)
func handleCheck(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}
	tokenStr := string(b)
	// Получаем токен в структурированном виде
	claims := &authlib.AuthClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) { return &authPrivateKey.PublicKey, nil })

	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	if !token.Valid {
		httptool.SetStatus401(w, "Invalid token")
		return
	}

	resp := &CheckResponse{
		UserUuid: claims.UserUuid,
		UserRole: claims.UserRole,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		httptool.SetStatus500(w, err)
	}
}
