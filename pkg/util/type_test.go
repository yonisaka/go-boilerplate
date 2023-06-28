// Package util
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSameType(t *testing.T) {
	assert.Equal(t, true, IsSameType(struct{}{}, struct{}{}))
	assert.Equal(t, false, IsSameType(struct{}{}, int64(1)))
}
