package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var eventSchemaRepo map[EventSchemaKey]string

func init() {
	eventSchemaRepo = make(map[EventSchemaKey]string)
}

// Загружаем каталог схем
// Для простоты - добавление новой схемы выполняется вручную в файл repository.json
// После добавления требуется перезапуск сервиса
func initSchemaRepo() {
	repoFileName := fmt.Sprintf("%s/repository.json", schemaRepositoryPath)

	source, err := os.ReadFile(repoFileName)
	if err != nil {
		log.Fatalln(err)
	}

	var items []EventRepoItem
	err = json.Unmarshal(source, &items)
	if err != nil {
		log.Fatalln(err)
	}

	for _, item := range items {
		eventSchemaRepo[EventSchemaKey{
			EventType:    item.EventType,
			EventVersion: item.EventVersion,
		}] = item.SchemaPath
	}
}
