package main

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestRunCmd(t *testing.T) {
	t.Run("check succes", func(t *testing.T) {
		result := RunCmd([]string{"echo", "1"}, map[string]EnvValue{"TEST": {Value: "FOOBAR", NeedRemove: false}})
		require.NotEqual(t, 1, result)
		require.Equal(t, 0, result)
	})
}
