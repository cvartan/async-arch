package main

import (
	model "async-arch/internal/domain/analisys"
	authmodel "async-arch/internal/domain/auth"
	"async-arch/internal/lib/auth"
	"async-arch/internal/lib/base"
	"async-arch/internal/lib/httptool"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func initHandlers() {
	base.App.HandleFunc("GET /api/v1/users", auth.WithAuth(handleGetUserStatRequest, []authmodel.UserRole{"ADMIN"}))
	base.App.HandleFunc("GET /api/v1/tasks", auth.WithAuth(handleGetTaskStatRequest, []authmodel.UserRole{"ADMIN"}))
}

const GetManagementFeeQueryTemplate = `
select 
	(select sum(t.value) from analisys."transaction" t where t."type" = 'DEBITING' and date_trunc('day',t."time")=current_date) debiting_sum, 
	(select sum(t.value) from analisys."transaction" t where t."type" = 'VALUE' and date_trunc('day',t."time")=current_date) value_sum 
`
const getUserListQueryTemplate = `
select u."uuid" ,u."name" from analisys."user" u
`
const getUserTransactionSumQueryTemplate = `
select 
	(select sum(t.value) from analisys."transaction" t where t."type" = 'DEBITING' and date_trunc('day',t."time")=current_date and t.user_uuid = ?) debiting_sum, 
	(select sum(t.value) from analisys."transaction" t where t."type" = 'VALUE' and date_trunc('day',t."time")=current_date and t.user_uuid = ?) value_sum 
`

// Получение отчета по зарабтокам менеджмента и попугам ушедших в минус
func handleGetUserStatRequest(w http.ResponseWriter, r *http.Request) {
	repo, _ := base.App.GetDomainRepository("analisys")

	result, err := repo.RawQuery(getUserListQueryTemplate)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}

	rows := result.(*sql.Rows)

	var users []*model.User

	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.Uuid, &user.Name)
		if err != nil {
			log.Fatalln(err)
		}
		users = append(users, user)
	}
	rows.Close()

	var lazyUsers []*LazyUserInfo
	var lazyUserCount int

	for _, user := range users {
		result, err := repo.RawQuery(getUserTransactionSumQueryTemplate, user.Uuid, user.Uuid)

		var debSum, valSum int
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}
		rows := result.(*sql.Rows)
		for rows.Next() {
			var debSumNullable, valSumNullable sql.Null[int]
			err := rows.Scan(&debSumNullable, &valSumNullable)
			if err != nil {
				log.Fatalln(err)
			}
			if debSumNullable.Valid {
				debSum = debSumNullable.V
			}

			if valSumNullable.Valid {
				valSum = valSumNullable.V
			}
		}
		rows.Close()

		balance := debSum - valSum
		if balance > 0 {
			lazyUser := &LazyUserInfo{
				Uuid:    user.Uuid,
				Name:    user.Name,
				Penalty: balance,
			}

			lazyUsers = append(lazyUsers, lazyUser)
			lazyUserCount++
		}
	}

	result, err = repo.RawQuery(GetManagementFeeQueryTemplate)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}
	rows = result.(*sql.Rows)

	var debSum, valSum int
	for rows.Next() {
		var debSumNullable, valSumNullable sql.Null[int]
		err := rows.Scan(&debSumNullable, &valSumNullable)
		if err != nil {
			log.Fatalln()
		}
		if debSumNullable.Valid {
			debSum = debSumNullable.V
		}
		if valSumNullable.Valid {
			valSum = valSumNullable.V
		}
	}
	rows.Close()

	response := &GetStatByUserResponse{
		ManagementFee:  debSum - valSum,
		LazyUsers:      lazyUsers,
		LazyUsersCount: lazyUserCount,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		httptool.SetStatus500(w, err)
	}
}

const MaxPriceByDaysQueryTemplate = `
select date_trunc('day',t.complete_time) "date",max(t.completed_price) 
from analisys.task t 
where t.is_complete = true and t.complete_time between ? and ?
group by date_trunc('day',t.complete_time)
`
const MaxPriceByPeriodQueryTemplate = `
select max(t.completed_price) 
from analisys.task t 
where t.is_complete = true and t.complete_time between ? and ?
`

// Получение отчета по самым дорогим завершенным задачам
func handleGetTaskStatRequest(w http.ResponseWriter, r *http.Request) {
	periodDates := r.URL.Query()
	dateFrom := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	dateFromQueryValue := periodDates.Get("from")
	if dateFromQueryValue != "" {
		var err error
		dateFrom, err = time.Parse(time.DateOnly, dateFromQueryValue)
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}
	}
	dateTo := time.Now()
	dateToQueryValue := periodDates.Get("to")
	if dateToQueryValue != "" {
		var err error
		dateTo, err = time.Parse(time.DateOnly, dateToQueryValue)
		if err != nil {
			httptool.SetStatus500(w, err)
			return
		}
		dateTo = dateTo.Add(time.Hour * 24)
	}

	repo, _ := base.App.GetDomainRepository("analisys")
	result, err := repo.RawQuery(MaxPriceByDaysQueryTemplate, dateFrom, dateTo)
	if err != nil {
		httptool.SetStatus500(w, err)
	}

	rows := result.(*sql.Rows)

	var pricesByDay []*DayMaxPrice

	for rows.Next() {
		priceByDay := &DayMaxPrice{}
		err := rows.Scan(&priceByDay.Date, &priceByDay.MaxPrice)
		if err != nil {
			log.Fatalln(err)
		}
		pricesByDay = append(pricesByDay, priceByDay)
	}

	rows.Close()

	maxPrice := int(0)

	result, err = repo.RawQuery(MaxPriceByPeriodQueryTemplate, dateFrom, dateTo)
	if err != nil {
		httptool.SetStatus500(w, err)
	}

	rows = result.(*sql.Rows)

	var maxPriceNullable sql.Null[int]
	for rows.Next() {
		err := rows.Scan(&maxPriceNullable)
		if err != nil {
			log.Fatalln(err)
		}
		if maxPriceNullable.Valid {
			maxPrice = maxPriceNullable.V
		}
	}
	rows.Close()

	response := &GetStatByTasksResponse{
		MaxPricesByDate: pricesByDay,
		TotalMaxPrice:   maxPrice,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalln(err)
	}
}
