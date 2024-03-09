package main

import "time"

type GetStatByUserResponse struct {
	ManagementFee  int             `json:"managementFee"`
	LazyUsersCount int             `json:"lazyUsersCount"`
	LazyUsers      []*LazyUserInfo `json:"lazyUsers"`
}

type LazyUserInfo struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Penalty int    `json:"penalty"`
}

type GetStatByTasksResponse struct {
	TotalMaxPrice   int            `json:"maxPriceByPeriod"`
	MaxPricesByDate []*DayMaxPrice `json:"maxPricesByDate"`
}

type DayMaxPrice struct {
	Date     time.Time `json:"date"`
	MaxPrice int       `json:"maxPrice"`
}
