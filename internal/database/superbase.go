package database

import (
	"log"
	"os"

	"github.com/nedpals/supabase-go"
)

var SupabaseClient *supabase.Client

func InitSupabase() {

	url := os.Getenv("DATABASE_URL")

	if url == "" {
		log.Fatal("DATABASE_URL fehlt in .env")
	}

	SupabaseClient = supabase.CreateClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
	)

	log.Println("Supabase verbunden")
}
