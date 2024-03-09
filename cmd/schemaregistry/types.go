package main

type EventSchemaKey struct {
	EventType    string
	EventVersion string
}

type EventRepoItem struct {
	EventType    string `json:"event"`
	EventVersion string `json:"version"`
	SchemaPath   string `json:"path"`
}

type EventSchemaRequest struct {
	EventType    string `json:"event"`
	EventVersion string `json:"version"`
}
