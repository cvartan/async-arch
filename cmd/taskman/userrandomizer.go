package main

import (
	base "async-arch/internal/lib/base"
	"database/sql"
	"log"
	"math/rand"
)

// UserRandomizer - определения случайного пользователя для назначения в задачу
type UserRandomizer struct {
	users []string
}

// Создание списка пользователей, которых можно назначить задаче
func CreateUserRandomizer() *UserRandomizer {
	var uuids []string
	repo, _ := base.App.GetDomainRepository("task")
	// Получение списка пользователей из БД
	result, err := repo.RawQuery("select uuid from task.user where role not in ('ADMIN','MANAGER')")
	if err != nil {
		log.Fatal(err)
	}
	rows := result.(*sql.Rows)
	defer rows.Close()
	for rows.Next() {
		var uuid string
		err = rows.Scan(&uuid)
		if err != nil {
			log.Fatal(err)
		}
		uuids = append(uuids, uuid)
	}

	return &UserRandomizer{
		users: uuids,
	}
}

// Получение идентификатора случайного пользователя
func (r *UserRandomizer) Uuid() string {
	return r.users[rand.Intn(len(r.users)-1)]
}
