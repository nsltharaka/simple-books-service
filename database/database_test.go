package database

import (
	"os"
	"testing"
)

func TestConnection(t *testing.T) {

	t.Run("returns error if SQLITE_FILENAME is not set", func(t *testing.T) {
		os.Unsetenv("SQLITE_FILENAME")
		_, err := Connection()
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("returns db instance if SQLITE_FILENAME is set", func(t *testing.T) {
		tempFile := "test_temp.db"
		os.Setenv("SQLITE_FILENAME", tempFile)
		defer os.Remove(tempFile)

		db, err := Connection()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if db == nil {
			t.Fatal("expected db instance, got nil")
		}
	})
}
