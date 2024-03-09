package main

import (
	"async-arch/internal/lib/base"
	"async-arch/internal/lib/sysenv"
	"log"
)

var schemaRepositoryPath string = sysenv.GetEnvValue(
	"JSON_SCHEMA_REPO+PATH",
	"D:/Projects/async-arch/api/event",
)

func main() {
	initSchemaRepo()

	err := base.App.InitHTTPServer("", 8094)
	if err != nil {
		log.Fatalln(err)
	}

	base.App.HandleFunc("GET /api/v1/schema/{event}/{version}", handleGetSchema)

	base.App.Hold()
}
