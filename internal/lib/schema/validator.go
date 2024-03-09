package schema

import (
	"async-arch/internal/lib/sysenv"
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

type SchemaValidator struct {
	schemaService     *http.Client
	schemaServiceAddr string
	schemaCache       map[EventSchemaKey]*jsonschema.Schema
}

func CreateSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemaService: &http.Client{},
		schemaServiceAddr: fmt.Sprintf(
			"%s://%s",
			sysenv.GetEnvValue("SCHEMA_SERVICE_METHOD", "http"),
			sysenv.GetEnvValue("SCHEMA_SERVICE_ADDR", "localhost:8094"),
		),
		schemaCache: make(map[EventSchemaKey]*jsonschema.Schema),
	}
}

func (v *SchemaValidator) Validate(eventType, eventVersion, jsonObject string) error {
	var (
		schema *jsonschema.Schema
		ok     bool
	)

	schema, ok = v.schemaCache[EventSchemaKey{
		EventType:    eventType,
		EventVersion: eventVersion,
	}]
	if !ok {
		// Если схему не нашли, то добавляем запрашиваем ее в сервисе ServiceRegistry
		req, err := http.NewRequest("GET", fmt.Sprintf(
			"%s/api/v1/schema/%s/%s",
			v.schemaServiceAddr,
			eventType,
			eventVersion,
		), nil)
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
	return schema.Validate(object)
}
