package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/fawwns/antibruteforce/internal/config"
	"github.com/fawwns/antibruteforce/internal/httpapi"
	"github.com/fawwns/antibruteforce/internal/postgres"
	"github.com/fawwns/antibruteforce/internal/service"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

var client = &http.Client{}

type CheckResponse struct {
	OK bool `json:"ok"`
}

func startTestServer() *http.Server {
	cfg := &config.Config{
		Port:          "8080",
		LimitLogin:    1000,
		LimitPassword: 1000,
		LimitIP:       3,

		Postgres: config.PostgresConfig{
			Host:     "127.0.0.1",
			Port:     "5432",
			User:     "antibruteforce",
			Password: "secret",
			DBName:   "antibruteforce",
		},
	}

	pg, err := postgres.NewPostgres(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		panic(err)
	}

	if err := pg.Migrate(); err != nil {
		panic(err)
	}

	svc := service.New(cfg, pg)
	h := httpapi.NewHandler(svc)
	router := httpapi.NewRouter(h)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		_ = srv.ListenAndServe()
	}()

	time.Sleep(200 * time.Millisecond)

	return srv
}

func doCheck(t *testing.T, login, password, ip string) bool {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/check?login="+login+"&password="+password+"&ip="+ip,
		nil,
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var body CheckResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))

	return body.OK
}

func post(t *testing.T, path string) {
	t.Helper()

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+path,
		bytes.NewBuffer(nil),
	)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}
