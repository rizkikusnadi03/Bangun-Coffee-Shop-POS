package database

import (
	"strings"
	"testing"
)

func TestNewMySQLConnectionFailure(t *testing.T) {
	// DSN ke host/port yang tidak ada untuk menguji error handling pada db.Ping()
	invalidDSN := "invalid_user:invalid_password@tcp(127.0.0.1:54321)/non_existent_db?timeout=1s"

	db, err := NewMySQL(invalidDSN)
	if err == nil {
		if db != nil {
			db.Close()
		}
		t.Fatal("expected error when connecting to invalid mysql server, got nil")
	}

	if !strings.Contains(err.Error(), "failed to ping mysql database") {
		t.Errorf("expected descriptive ping failure error, got: %v", err)
	}
}
