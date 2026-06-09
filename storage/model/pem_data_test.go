package model

import (
	"bytes"
	"testing"
)

func TestPEMDataColumnType(t *testing.T) {
	tests := []struct {
		dialect string
		want    string
	}{
		{dialect: "mysql", want: "longblob"},
		{dialect: "postgres", want: "bytea"},
		{dialect: "sqlite", want: "blob"},
		{dialect: "other", want: "blob"},
	}

	for _, tt := range tests {
		if got := pemDataColumnType(tt.dialect); got != tt.want {
			t.Fatalf("pemDataColumnType(%q) = %q, want %q", tt.dialect, got, tt.want)
		}
	}
}

func TestPEMDataScanAndValue(t *testing.T) {
	var p PEMData

	if err := p.Scan([]byte("abc")); err != nil {
		t.Fatalf("Scan([]byte) returned error: %v", err)
	}
	if !bytes.Equal([]byte(p), []byte("abc")) {
		t.Fatalf("Scan([]byte) = %q, want %q", []byte(p), []byte("abc"))
	}

	val, err := p.Value()
	if err != nil {
		t.Fatalf("Value() returned error: %v", err)
	}
	if got, ok := val.([]byte); !ok || !bytes.Equal(got, []byte("abc")) {
		t.Fatalf("Value() = %#v, want []byte(%q)", val, "abc")
	}

	if err := p.Scan("xyz"); err != nil {
		t.Fatalf("Scan(string) returned error: %v", err)
	}
	if !bytes.Equal([]byte(p), []byte("xyz")) {
		t.Fatalf("Scan(string) = %q, want %q", []byte(p), []byte("xyz"))
	}

	if err := p.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) returned error: %v", err)
	}
	if p != nil {
		t.Fatalf("Scan(nil) = %#v, want nil", []byte(p))
	}
}

