package event_test

import (
	event "async-arch/internal/lib/event"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

type EmbededData struct {
	Name   string    `json:"name"`
	Weight int       `json:"weight"`
	Flash  time.Time `json:"flash"`
}

type StringWriter struct {
	source []byte
}

func (w *StringWriter) Write(p []byte) (int, error) {
	w.source = p
	return len(p), nil
}

func (w *StringWriter) String() string {
	return string(w.source)
}

func TestEmbededEvent(t *testing.T) {
	emb := EmbededData{
		Name:   "Test",
		Weight: 100,
		Flash:  time.Now(),
	}

	e := event.Event{
		EventID:   uuid.NewString(),
		EventType: "Test",
		Subject:   reflect.TypeOf(emb).String(),
		Sender:    "test",
		CreatedAt: time.Now(),
	}

	msg := event.EventMessage{
		Event: e,
		Data:  emb,
	}

	w := StringWriter{}
	err := json.NewEncoder(&w).Encode(msg)
	if err != nil {
		t.Fatal(err)
	}

	s := w.String()

	var ne event.EventMessage
	err = json.NewDecoder(strings.NewReader(s)).Decode(&ne)
	if err != nil {
		log.Fatal(err)
	}

}
