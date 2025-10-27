package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("directory is not exists", func(t *testing.T) {
		_, err := ReadDir("dir_is_not_exists")
		require.Error(t, err)
	})

	t.Run("empty directory", func(t *testing.T) {
		dir := "testdata/empty"

		err := os.Mkdir(dir, 0o777)
		require.Nil(t, err)
		defer func() {
			err = os.Remove(dir)
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, Environment{}, env)
	})

	t.Run("success case", func(t *testing.T) {
		expected := Environment{
			"BAR":   {"bar", false},
			"UNSET": {"", true},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
		}

		actual, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("incorrect filename", func(t *testing.T) {
		dir := "testdata/invalid_env"

		err := os.Mkdir(dir, os.FileMode(0o755))
		require.NoError(t, err)
		defer func() {
			err = os.RemoveAll(dir)
			require.NoError(t, err)
		}()

		f, err := os.Create(path.Join(dir, "NAME=INVALID"))
		require.Nil(t, err)
		defer func() {
			err = f.Close()
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)
		require.Zero(t, env)
		require.Error(t, err)
	})
}
