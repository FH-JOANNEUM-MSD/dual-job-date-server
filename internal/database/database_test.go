package database

import "testing"

func TestSupabaseClient_BeforeInit(t *testing.T) {
	previous := SupabaseClient
	SupabaseClient = nil
	defer func() {
		SupabaseClient = previous
	}()

	if SupabaseClient != nil {
		t.Fatalf("expected SupabaseClient to be nil before init")
	}
}
