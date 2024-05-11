package main

import (
	authmodel "async-arch/internal/domain/auth"
	auth "async-arch/internal/lib/auth"
	base "async-arch/internal/lib/base"
	"async-arch/internal/lib/httputils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func initHandlers() {
	base.App.HandleFunc("GET /api/v1/log", auth.WithAuth(handleGetLog, nil))
	base.App.HandleFunc("GET /api/v1/stat", auth.WithAuth(handleGetStat, []authmodel.UserRole{authmodel.ACCOUNTER, authmodel.ADMIN}))
}

const getBalanceQueryTemplate = `
select u.balance 
from accounting."user" u 
where u."uuid" = ?
`

const getLogQueryTemplate = `
select t."uuid" ,t."time" ,t."task_uuid", t."type" ,t.value 
from accounting."transaction" t  
where t.user_uuid = ? and date_trunc('day',t."time") = current_date 
order by t."time"
`

// Получение баланса и лога операций для пользователя
func handleGetLog(w http.ResponseWriter, r *http.Request) {

	userUuid := r.Header.Get("X-Auth-User-UUID")

	repo, _ := base.App.GetDomainRepository("accounting")
	result, err := repo.RawQuery(getBalanceQueryTemplate, userUuid)

	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		log.Fatalln("result is not rows")
	}

	var balance int
	for rows.Next() {
		err := rows.Scan(&balance)
		if err != nil {
			log.Fatalln(err)
		}
	}
	rows.Close()

	result, err = repo.RawQuery(getLogQueryTemplate, userUuid)
	if err != nil {
		httputils.SetStatus500(w, err)
	}

	rows, ok = result.(*sql.Rows)
	if !ok {
		log.Fatalln("result is not rows")
	}

	var items []TransactionLogItem

	for rows.Next() {
		item := &TransactionLogItem{}
		err := rows.Scan(&item.Uuid, &item.Time, &item.TaskUuid, &item.Type, &item.Value)
		if err != nil {
			log.Fatalln(err)
		}
		items = append(items, *item)
	}
	rows.Close()

	response := &GetLogResponse{
		Balance: balance,
		Log:     items,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

}

const getStatQueryTemplate = `
select 
	(select sum(t.value) from accounting."transaction" t where t."type" = 'DEBITING' and date_trunc('day',t."time")=current_date) debiting_sum, 
	(select sum(t.value) from accounting."transaction" t where t."type" = 'VALUE' and date_trunc('day',t."time")=current_date) value_sum 
`

// Получение суммы заработанной менеджментом
func handleGetStat(w http.ResponseWriter, h *http.Request) {
	repo, _ := base.App.GetDomainRepository("accounting")

	result, err := repo.RawQuery(getStatQueryTemplate)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		log.Fatalln("result is not rows")
	}

	debSum, valSum := int(0), int(0)

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

	w.Write([]byte(strconv.Itoa(debSum - valSum)))

}
