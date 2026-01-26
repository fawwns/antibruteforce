package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/fawwns/antibruteforce/internal/logger"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.Info = log.New(io.Discard, "", 0)
	logger.Warn = log.New(io.Discard, "", 0)
	logger.Error = log.New(io.Discard, "", 0)
}

func TestRateLimitByLogin(t *testing.T) {
	login := "login-limit"
	password := "pass"
	ip := "10.0.0.1"

	blocked := false

	for i := 0; i < 50; i++ {
		ok := doCheck(t, login, password, ip)
		if !ok {
			blocked = true
			break
		}
	}

	require.True(t, blocked, "service must eventually block by login")
}

func TestRateLimitByPassword(t *testing.T) {
	password := "same-pass"

	blocked := false

	for i := 0; i < 50; i++ {
		ok := doCheck(
			t,
			"user"+string(rune('a'+i)),
			password,
			"10.0.0."+string(rune('1'+i)),
		)
		if !ok {
			blocked = true
			break
		}
	}

	require.True(t, blocked, "service must eventually block by password")
}

func TestRateLimitByIP(t *testing.T) {
	srv := startTestServer()
	defer srv.Close()

	ip := "10.0.0.100"
	blocked := false

	for i := 0; i < 1001; i++ {
		login := fmt.Sprintf("user-%d", i)
		password := fmt.Sprintf("pass-%d", i)

		resp, err := http.Get(
			fmt.Sprintf("%s/check?login=%s&password=%s&ip=%s",
				baseURL, login, password, ip,
			),
		)
		require.NoError(t, err)

		var r struct {
			OK bool `json:"ok"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&r))
		resp.Body.Close()

		if !r.OK {
			blocked = true
			break
		}
	}

	require.True(t, blocked, "service must eventually block by IP")
}
