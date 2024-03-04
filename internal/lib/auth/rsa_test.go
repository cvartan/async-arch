package auth_test

import (
	auth "async-arch/internal/lib/auth"
	"testing"
)

func TestRSARandom(t *testing.T) {
	_, err := auth.CreateJwtTokenChecker("http", "localhost:8090", "GET", "/api/v1/key")
	if err != nil {
		t.Fatal(err)
	}
}
