package main

import "time"

// Ответ с отчетом по заработкам менеджмента и "ленивым" попугам
type GetStatByUserResponse struct {
	ManagementFee  int             `json:"managementFee"`
	LazyUsersCount int             `json:"lazyUsersCount"`
	LazyUsers      []*LazyUserInfo `json:"lazyUsers"`
}

// Информация по "ленивому" попугу (ушедшему в минус)
type LazyUserInfo struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Penalty int    `json:"penalty"`
}

// Ответ с отчетом по дорогим закрытым задачам
type GetStatByTasksResponse struct {
	TotalMaxPrice   int            `json:"maxPriceByPeriod"`
	MaxPricesByDate []*DayMaxPrice `json:"maxPricesByDate"`
}

type DayMaxPrice struct {
	Date     time.Time `json:"date"`
	MaxPrice int       `json:"maxPrice"`
}
