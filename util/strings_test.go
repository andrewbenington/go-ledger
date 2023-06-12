package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDateFromTitle(t *testing.T) {
	m, d := ExtractDateFromTitle("Transaction 06/02")
	assert.Equal(t, m, 6)
	assert.Equal(t, d, 2)
}
