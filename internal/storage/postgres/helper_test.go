package postgres

import "testing"

func TestRebind(t *testing.T) {
	got := rebind(`?, ?, ?, ?`)
	expected := "$1, $2, $3, $4"
	if got != expected {
		t.Errorf("expected %s, but got %s", expected, got)
	}

	got = rebind(`SELECT ?, ?, ?, ?`)
	expected = "SELECT $1, $2, $3, $4"
	if got != expected {
		t.Errorf("expected %s, but got %s", expected, got)
	}
}
