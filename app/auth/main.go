package main

import (
	"async-arch/lib/base"
	"log"
)

func main() {
	if err := base.App.InitHTTPServer("", 8090); err != nil {
		log.Fatalln(err)
	}

	initModel()
	initHandlers()

	// Запускаем приложение
	base.App.Hold()
}
