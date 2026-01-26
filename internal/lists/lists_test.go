package lists_test

import (
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/fawwns/antibruteforce/internal/lists"
	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/postgres"
	"github.com/joho/godotenv"
)

func init() {
	logger.Info = log.New(io.Discard, "", 0)
	logger.Warn = log.New(io.Discard, "", 0)
	logger.Error = log.New(io.Discard, "", 0)
	_ = godotenv.Load("../../.env")
}

func setupPostgres(t *testing.T) *postgres.Postgres {
	host := os.Getenv("LOCAL_POSTGRES_HOST")
	port := os.Getenv("LOCAL_POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	pg, err := postgres.NewPostgres(host, port, user, password, dbname)
	if err != nil {
		t.Fatalf("failed to connect to Postgres: %v", err)
	}

	if err := pg.Migrate(); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return pg
}

func TestAddToWhitelist(t *testing.T) {

	pg := setupPostgres(t)

	l := lists.New(pg)

	err := l.AddToWhitelist("192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ip := net.ParseIP("192.168.1.10")
	if !l.InWhitelist(ip) {
		t.Fatalf("expected IP to be in whitelist")
	}
}

func TestAddToBlacklist(t *testing.T) {
	pg := setupPostgres(t)

	l := lists.New(pg)

	err := l.AddToBlacklist("192.168.1.0/24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ip := net.ParseIP("192.168.1.10")
	if !l.InBlacklist(ip) {
		t.Fatalf("expected IP to be in whitelist")
	}
}

func TestRemoveFromWhitelist(t *testing.T) {
	pg := setupPostgres(t)

	l := lists.New(pg)

	l.AddToWhitelist("192.168.0.0/16")
	l.RemoveFromWhitelist("192.168.0.0/16")

	ip := net.ParseIP("192.168.1.1")
	if l.InWhitelist(ip) {
		t.Fatalf("expected IP to NOT be in whitelist after removal")
	}
}

func TestRemoveFromBlacklist(t *testing.T) {
	pg := setupPostgres(t)

	l := lists.New(pg)

	l.AddToBlacklist("10.0.0.0/8")
	l.RemoveFromBlacklist("10.0.0.0/8")

	ip := net.ParseIP("10.0.0.1")
	if l.InBlacklist(ip) {
		t.Fatalf("expected IP to NOT be in blacklist after removal")
	}
}

func TestInWhitelistFalse(t *testing.T) {
	pg := setupPostgres(t)

	l := lists.New(pg)
	ip := net.ParseIP("1.2.3.4")

	if l.InWhitelist(ip) {
		t.Fatalf("expected false for empty whitelist")
	}
}

func TestInBlacklistFalse(t *testing.T) {
	pg := setupPostgres(t)

	l := lists.New(pg)
	ip := net.ParseIP("1.2.3.4")

	if l.InBlacklist(ip) {
		t.Fatalf("expected false for empty blacklist")
	}
}
