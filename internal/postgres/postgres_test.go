package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/fawwns/antibruteforce/internal/postgres"
	"github.com/joho/godotenv"
)

func TestPostgresStorage(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, relying on system environment")
	}

	host := os.Getenv("LOCAL_POSTGRES_HOST")
	port := os.Getenv("LOCAL_POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	fmt.Println("LOCAL_POSTGRES_HOST=", os.Getenv("LOCAL_POSTGRES_HOST"))
	fmt.Println("LOCAL_POSTGRES_PORT=", os.Getenv("LOCAL_POSTGRES_PORT"))
	fmt.Println("POSTGRES_USER=", os.Getenv("POSTGRES_USER"))
	fmt.Println("POSTGRES_PASSWORD=", os.Getenv("POSTGRES_PASSWORD"))
	fmt.Println("POSTGRES_DB=", os.Getenv("POSTGRES_DB"))
	t.Logf("Connecting to Postgres: %s:%s user=%s db=%s", host, port, user, dbname)

	store, err := postgres.NewPostgres(host, port, user, password, dbname)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	if err := store.Migrate(); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	ctx := context.Background()
	cidr := "192.168.1.0/24"

	// Whitelist.
	if err := store.AddToWhitelist(ctx, cidr); err != nil {
		t.Fatalf("AddToWhitelist error: %v", err)
	}
	wlist, _ := store.GetWhitelist(ctx)
	found := false
	for _, c := range wlist {
		if c == cidr {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("cidr not found in whitelist")
	}

	if err := store.RemoveFromWhitelist(ctx, cidr); err != nil {
		t.Fatalf("RemoveFromWhitelist error: %v", err)
	}

	// Blacklist.
	if err := store.AddToBlacklist(ctx, cidr); err != nil {
		t.Fatalf("AddToBlacklist error: %v", err)
	}
	blist, _ := store.GetBlacklist(ctx)
	found = false
	for _, c := range blist {
		if c == cidr {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("cidr not found in blacklist")
	}

	if err := store.RemoveFromBlacklist(ctx, cidr); err != nil {
		t.Fatalf("RemoveFromBlacklist error: %v", err)
	}
}
