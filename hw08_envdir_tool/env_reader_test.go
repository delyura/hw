package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestReadDir(t *testing.T) {
	t.Run("correct reading", func(t *testing.T) {
		expected := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		actual, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("no exist directory", func(t *testing.T) {
		_, err := ReadDir("no exist directory")
		require.Error(t, err)
	})

	t.Run("empty directory", func(t *testing.T) {
		dir := "testdata/foobar"
		err := os.Mkdir(dir, os.FileMode(0o755))
		require.Nil(t, err)
		defer func() {
			err = os.Remove(dir)
			require.NoError(t, err)
		}()

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, Environment{}, env)
	})

	t.Run("invalid name", func(t *testing.T) {
		dir := "testdata/invalid_name"
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
		require.Equal(t, os.ErrInvalid, err)
	})
}
