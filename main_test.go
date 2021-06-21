package main

import (
	"testing"
)

func TestGenerateDSNFromConfig(t *testing.T) {
	cfg := &ClientConfig{
		User:     "user",
		Pass:     "pass",
		Host:     "localhost",
		Port:     1234,
		Database: "db",
	}

	expected := "user:pass@tcp(localhost:1234)/db?multiStatements=true"
	dsn := cfg.DSN()

	if dsn != expected {
		t.Errorf("generated dsn unexpected, got: '%s' expected: '%s'", dsn, expected)
	}
}

func TestGenerateDSNWithTLSFromConfig(t *testing.T) {
	cfg := &ClientConfig{
		CA:         "123",
		ClientKey:  "123",
		ClientCert: "123",
		User:       "user",
		Pass:       "pass",
		Host:       "localhost",
		Port:       1234,
		Database:   "db",
	}

	expected := "user:pass@tcp(localhost:1234)/db?multiStatements=true&tls=custom"
	dsn := cfg.DSN()

	if dsn != expected {
		t.Errorf("generated dsn unexpected, got: '%s' expected: '%s'", dsn, expected)
	}
}
