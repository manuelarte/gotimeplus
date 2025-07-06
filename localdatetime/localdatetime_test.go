package localdatetime

import (
	"testing"

	"github.com/manuelarte/gotimeplus/localdate"
	"github.com/manuelarte/gotimeplus/localtime"
)

func TestLocalDateTime_Equal(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		localDateTime LocalDateTime
		other         LocalDateTime
		expected      bool
	}{
		"same local date time": {
			localDateTime: New(
				localdate.New(2000, 1, 1),
				localtime.New(0, 0, 0, 0),
			),
			other: New(
				localdate.New(2000, 1, 1),
				localtime.New(0, 0, 0, 0),
			),
			expected: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := test.localDateTime.Equal(test.other)

			if test.expected != actual {
				t.Errorf("\nExpected: %v\nActual: %v", test.expected, actual)
			}
		})
	}
}
