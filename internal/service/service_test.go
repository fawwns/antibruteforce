package service_test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/postgres"
	"github.com/fawwns/antibruteforce/internal/service"
	"github.com/joho/godotenv"
)

func init() {
	logger.Warn = log.New(io.Discard, "", 0)
	logger.Info = log.New(io.Discard, "", 0)
	logger.Error = log.New(io.Discard, "", 0)
}

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

func TestServiceCheck(t *testing.T) {

	cfg := &config.Config{
		LimitLogin:    2,
		LimitPassword: 2,
		LimitIP:       2,
	}

	pg := setupPostgres(t)
	srv := service.New(cfg, pg)

	ip := "1.2.3.4"
	login := "user"
	pass := "pwd"

	for i := 0; i < 2; i++ {
		if !srv.Check(login, pass, ip) {
			t.Fatalf("expected Check() = true on iteration %d", i)
		}
	}

	if srv.Check(login, pass, ip) {
		t.Fatalf("expected Check() = false on iteration")
	}
}

func TestServiceWhitelist(t *testing.T) {
	cfg := &config.Config{
		LimitLogin:    1,
		LimitPassword: 1,
		LimitIP:       1,
	}
	pg := setupPostgres(t)
	srv := service.New(cfg, pg)

	srv.AddToWhitelist("192.168.0.0/24")

	ok := srv.Check("u", "p", "192.168.0.10")
	if !ok {
		t.Fatalf("expected whitelisted IP to bypass limiter")
	}
}

func TestServiceBlacklist(t *testing.T) {
	pg := setupPostgres(t)

	svc := service.New(&config.Config{}, pg)

	svc.AddToBlacklist("10.0.0.0/8")

	ok := svc.Check("u", "p", "10.1.2.3")
	if ok {
		t.Fatalf("expected blacklisted IP to be rejected")
	}
}
