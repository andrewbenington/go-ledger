package ledger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLabel(t *testing.T) {
	t.Run("case insensitive", func(t *testing.T) {
		allLabels = []Label{{
			Name:     "test",
			Keywords: []string{"memo"},
		}}
		re, err := allLabels[0].RegExp()
		assert.NoError(t, err)
		allLabels[0].re = re
		label := FindLabel("LONG MEMO")
		assert.Equal(t, "test", label)
	})
}
