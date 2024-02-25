package main

import (
	"async-arch/lib/base"
	model "async-arch/model/domain/auth"
	"async-arch/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	ParrotBeak string `json:"beak"`
	Name       string `json:"name"`
	EMail      string `json:"email"`
	Role       string `json:"role"`
}

// HandleCreateUser - обработка запроса добавления пользователя
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	var userRq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRq); err != nil {
		util.GenerateHTTP500Response(w, err)
		return
	}

	user := model.CreateNewUser(userRq.ParrotBeak, userRq.Name, userRq.EMail, userRq.Role)

	result := db.Create(&user)
	if result.Error != nil {
		util.GenerateHTTP500Response(w, result.Error)
		return
	}

	// TODO: добавить отправку события добавления пользователя

}

// HandleAuthentificate - метод аутентификации попугая
func HandleAuthentificate(w http.ResponseWriter, r *http.Request) {

	parrotBeak := r.Header.Get("Authorization")

	var user model.User

	result := db.Where("parrot_beak=?", parrotBeak).Find(&user)
	if result.Error != nil {
		w.WriteHeader(401)
		return
	}

	userId := user.ID

	// Ищем открытую сессию для этого пользователя
	for id, params := range sessions {
		if params.UserID == userId && params.Expires.After(time.Now()) {
			params.Expires = time.Now().Add(time.Minute * 5)
			w.Header().Add("X-Session-ID", id)
			return
		}
	}

	// Если ничего не нашли - создаем новую сессию
	sessionId := uuid.NewString()
	sessions[sessionId] = sessionParams{
		UserID:  userId,
		Expires: time.Now().Add(time.Minute * 5),
	}
	// Возвращаем ключ сессии в заголовке X-Session-ID
	w.Header().Add("X-Session-ID", sessionId)
}

type AuthRequest struct {
	SessionId      string `json:"session"`
	PermissionCode string `json:"permission"`
}

// HandleAuthorize - метод авторизация пользователя
func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	var authReq AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	sessionParams, ok := sessions[authReq.SessionId]
	if !ok {
		w.WriteHeader(401)
		return
	}
	if sessionParams.Expires.Before(time.Now()) {
		w.WriteHeader(401)
		return
	}

	var user model.User
	result := db.First(&user, sessionParams.UserID)
	if result.Error != nil {
		w.WriteHeader(401)
		return
	}

	if roleSet, ok := permissions[authReq.PermissionCode]; ok {
		for _, value := range roleSet {
			if value == user.Role || value == "ALL" {
				w.WriteHeader(200)
				return
			}
		}
	}

	w.WriteHeader(401)
}

func initHandlers() {
	base.App.HandleFunc("POST /api/v1/users", HandleCreateUser)
	base.App.HandleFunc("POST /api/v1/login", HandleAuthentificate)
	base.App.HandleFunc("POST /api/v1/auth", HandleAuthorize)
}
