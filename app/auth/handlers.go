package main

import (
	"async-arch/lib/base"
	model "async-arch/model/domain/auth"
	"async-arch/util"
	"encoding/json"
	"log"
	"net/http"
)

type CreateUserRequest struct {
	ParrotBeak string   `json:"beak"`
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	Permisions []string `json:"permissions"`
}

// HandleCreateUser - обработка запроса добавления пользователя
func HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	var userRq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRq); err != nil {
		util.GenerateHTTP500Response(w, err)
	}

	user := model.CreateNewUser(userRq.ParrotBeak, userRq.Name, userRq.Role)

	for _, value := range userRq.Permisions {
		user.Permissions = append(user.Permissions, model.UserPermission{Permission: value})
	}

	db.Create(&user)

	log.Println(user.ID)

	// TODO: добавить отправку события добавления пользователя

}

func initHandlers() {
	base.App.HandleFunc("POST /api/v1/users/", HandleCreateUser)
}
