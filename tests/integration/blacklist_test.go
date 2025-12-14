package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlacklistAlwaysBlocks(t *testing.T) {
	post(t, "/blacklist/add?cidr=10.10.0.0/16")

	ok := doCheck(t, "any", "any", "10.10.5.5")
	require.False(t, ok, "blacklisted IP must be blocked")

	post(t, "/blacklist/remove?cidr=10.10.0.0/16")

	ok = doCheck(t, "any", "any", "10.10.5.5")
	require.True(t, ok, "IP should be allowed after removal")
}
