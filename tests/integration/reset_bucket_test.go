package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResetBucket(t *testing.T) {
	login := "reset-user"
	password := "reset-pass"
	ip := "172.16.0.1"

	blocked := false
	for i := 0; i < 100; i++ {
		ok := doCheck(t, login, password, ip)
		if !ok {
			blocked = true
			break
		}
	}
	require.True(t, blocked, "request must be blocked before reset")

	post(t, "/bucket/reset")

	ok := doCheck(t, login, password, ip)
	require.True(t, ok, "bucket should be reset")
}
