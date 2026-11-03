package redis

import (
	"strings"
	"testing"
)

func TestNewRedisConnectionFailure(t *testing.T) {
	// Koneksi ke host/port yang tidak ada untuk menguji error handling ping
	client, err := NewRedis("127.0.0.1:54321", "")
	if err == nil {
		if client != nil {
			_ = client.Close()
		}
		t.Fatal("expected error connecting to non-existent redis instance, got nil")
	}

	if !strings.Contains(err.Error(), "failed to ping redis server") {
		t.Errorf("expected descriptive ping failure message, got: %v", err)
	}
}
