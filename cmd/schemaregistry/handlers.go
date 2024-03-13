package main

import (
	"async-arch/internal/lib/httptool"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Получение схемы для события и версии события
// Файл со схемой расположен на диске
func handleGetSchema(w http.ResponseWriter, r *http.Request) {

	eventType := r.PathValue("event")
	eventVersion := r.PathValue("version")

	path, ok := eventSchemaRepo[EventSchemaKey{
		EventType:    eventType,
		EventVersion: eventVersion,
	}]
	if !ok {
		httptool.SetStatus500(w, errors.New("schema is not found"))
		return
	}

	schemaFile := fmt.Sprintf("%s/%s", schemaRepositoryPath, path)
	source, err := os.ReadFile(schemaFile)
	if err != nil {
		httptool.SetStatus500(w, err)
		return
	}
	w.Write(source)
}
