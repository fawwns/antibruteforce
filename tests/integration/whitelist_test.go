package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhitelistAlwaysAllows(t *testing.T) {
	post(t, "/whitelist/add?cidr=192.168.1.0/24")

	login := "wl-user"
	password := "wl-pass"
	ip := "192.168.1.10"

	for i := 0; i < 10; i++ {
		ok := doCheck(t, login, password, ip)
		require.True(t, ok, "whitelisted IP must always be allowed")
	}
}
