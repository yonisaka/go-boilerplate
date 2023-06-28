package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	t.Parallel()
	scenarios := []struct {
		Input    interface{}
		Expected string
	}{
		{
			Input:    int(1),
			Expected: "1",
		},
		{
			Input:    int8(2),
			Expected: "2",
		},
		{
			Input:    int16(3),
			Expected: "3",
		},
		{
			Input:    int32(4),
			Expected: "4",
		},
		{
			Input:    int64(5),
			Expected: "5",
		},
		{
			Input:    uint(6),
			Expected: "6",
		},
		{
			Input:    uint8(7),
			Expected: "7",
		},
		{
			Input:    uint16(8),
			Expected: "8",
		},
		{
			Input:    uint32(9),
			Expected: "9",
		},
		{
			Input:    uint64(10),
			Expected: "10",
		},
		{
			Input:    float32(11),
			Expected: "11",
		},
		{
			Input:    float64(12),
			Expected: "12",
		},
		{
			Input:    bool(true),
			Expected: "true",
		},
		{
			Input: struct {
			}{},
			Expected: "{}",
		},

		{
			Input:    string(`test`),
			Expected: "test",
		},
	}

	for _, sc := range scenarios {
		r := ToString(sc.Input)

		assert.Equal(t, sc.Expected, r)
	}
}
