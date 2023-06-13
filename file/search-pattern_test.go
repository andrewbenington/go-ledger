package file

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobDirectory(t *testing.T) {
	t.Run("Posix Home", func(t *testing.T) {
		t.Setenv("HOME", "/users/test")
		defer func() {
			runtimeOS = runtime.GOOS
		}()
		runtimeOS = "linux"
		dir := globDirectory("~/Documents")
		assert.Equal(t, "/users/test/Documents", dir)
	})
}
