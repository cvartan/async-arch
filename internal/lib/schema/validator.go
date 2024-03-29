package schema

import (
	ou "async-arch/internal/lib/osutils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

type EventSchemaKey struct {
	EventType    string
	EventVersion string
}

// Валидатор схем
type SchemaValidator struct {
	schemaService     *http.Client
	schemaServiceAddr string
	schemaCache       map[EventSchemaKey]*jsonschema.Schema
}

// Создание валидатора схем с нужными параметрами
func CreateSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemaService: &http.Client{},
		schemaServiceAddr: fmt.Sprintf(
			"%s://%s",
			ou.GetEnvValue("SCHEMA_SERVICE_METHOD", "http"),
			ou.GetEnvValue("SCHEMA_SERVICE_ADDR", "localhost:8094"),
		),
		schemaCache: make(map[EventSchemaKey]*jsonschema.Schema),
	}
}

// Валидация схемы события (данных события)
func (v *SchemaValidator) Validate(eventType, eventVersion, jsonObject string) error {
	var (
		schema *jsonschema.Schema
		ok     bool
	)

	// Ищем схему для этого события в кэше
	schema, ok = v.schemaCache[EventSchemaKey{
		EventType:    eventType,
		EventVersion: eventVersion,
	}]
	if !ok {
		// Если схему не нашли, то запрашиваем ее в сервисе ServiceRegistry
		req, err := http.NewRequest("GET", fmt.Sprintf(
			"%s/api/v1/schema/%s/%s",
			v.schemaServiceAddr,
			eventType,
			eventVersion,
		), nil)

		// Если схему не удалось получить то возвращаем ошибку
		if err != nil {
			return err
		}
		resp, err := v.schemaService.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 500 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return errors.New(string(body))
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		schema, err = jsonschema.CompileString("schema.json", string(body))
		if err != nil {
			return err
		}
		// Добавляем в кэш для последующего использования
		v.schemaCache[EventSchemaKey{
			EventType:    eventType,
			EventVersion: eventVersion,
		}] = schema
	}
	object := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonObject), &object)
	if err != nil {
		return err
	}
	// Вовзращаем результат валидации
	return schema.Validate(object)
}
