package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("incorrect cmd", func(t *testing.T) {
		code := RunCmd(make([]string, 0), make(Environment))
		require.Equal(t, 1, code)
	})
}
