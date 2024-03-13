package main

// Ключ схемы данных
type EventSchemaKey struct {
	EventType    string `json:"event"`
	EventVersion string `json:"version"`
}

// Запись о схеме
type EventRepoItem struct {
	EventType    string `json:"event"`
	EventVersion string `json:"version"`
	SchemaPath   string `json:"path"`
}

// Запрос схемы
type EventSchemaRequest struct {
	EventSchemaKey
}
