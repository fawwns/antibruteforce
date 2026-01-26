package httpapi_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/httpapi"
	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/fawwns/antibruteforce/internal/postgres"
	"github.com/fawwns/antibruteforce/internal/service"
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

func makeService(t *testing.T) *service.AntiBruteForceService {
	cfg := &config.Config{
		LimitIP:       10,
		LimitLogin:    10,
		LimitPassword: 10,
	}
	pg := setupPostgres(t)
	return service.New(cfg, pg)
}

func TestCheckHandler(t *testing.T) {
	svc := makeService(t)
	h := httpapi.NewHandler(svc)

	router := http.NewServeMux()
	router.HandleFunc("/check", h.Check)

	req := httptest.NewRequest(http.MethodGet, "/check?login=a&password=b&ip=1.2.3.4", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp httpapi.CheckResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("json: %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected OK=true")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/check?login=a&password=b", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr2.Code)
	}
}

func TestWhitelistHandlers(t *testing.T) {
	svc := makeService(t)
	h := httpapi.NewHandler(svc)

	router := http.NewServeMux()
	router.HandleFunc("/whitelist/add", h.AddWhitelist)
	router.HandleFunc("/whitelist/remove", h.RemoveFromWhitelist)

	reqAdd := httptest.NewRequest(http.MethodPost, "/whitelist/add?cidr=192.168.1.0/24", nil)
	rrAdd := httptest.NewRecorder()
	router.ServeHTTP(rrAdd, reqAdd)
	if rrAdd.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rrAdd.Code)
	}

	reqAddFail := httptest.NewRequest(http.MethodPost, "/whitelist/add", nil)
	rrAddFail := httptest.NewRecorder()
	router.ServeHTTP(rrAddFail, reqAddFail)
	if rrAddFail.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rrAddFail.Code)
	}

	reqRemove := httptest.NewRequest(http.MethodPost, "/whitelist/remove?cidr=192.168.1.0/24", nil)
	rrRemove := httptest.NewRecorder()
	router.ServeHTTP(rrRemove, reqRemove)
	if rrRemove.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rrRemove.Code)
	}
}

func TestBlacklistHandlers(t *testing.T) {
	svc := makeService(t)
	h := httpapi.NewHandler(svc)

	router := http.NewServeMux()
	router.HandleFunc("/blacklist/add", h.AddBlacklist)
	router.HandleFunc("/blacklist/remove", h.RemoveFromBlacklist)

	reqAdd := httptest.NewRequest(http.MethodPost, "/blacklist/add?cidr=10.0.0.0/8", nil)
	rrAdd := httptest.NewRecorder()
	router.ServeHTTP(rrAdd, reqAdd)
	if rrAdd.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rrAdd.Code)
	}

	reqRemove := httptest.NewRequest(http.MethodPost, "/blacklist/remove?cidr=10.0.0.0/8", nil)
	rrRemove := httptest.NewRecorder()
	router.ServeHTTP(rrRemove, reqRemove)
	if rrRemove.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rrRemove.Code)
	}
}

func TestResetBucketHandler(t *testing.T) {
	svc := makeService(t)
	h := httpapi.NewHandler(svc)

	router := http.NewServeMux()
	router.HandleFunc("/bucket/reset", h.ResetBucket)

	req := httptest.NewRequest(http.MethodPost, "/bucket/reset", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
