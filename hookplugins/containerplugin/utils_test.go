package containerplugin

import (
	"reflect"
	"sort"
	"testing"
)

func TestUniqueStringSlice(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{},
			expected: []string{},
		}, {
			input:    []string{"1"},
			expected: []string{"1"},
		}, {
			input:    []string{"3", "1", "2", "1", "2"},
			expected: []string{"1", "2", "3"},
		}, {
			input:    []string{"3", "1", "2"},
			expected: []string{"1", "2", "3"},
		},
	}

	for _, tc := range cases {
		got := UniqueStringSlice(tc.input)
		sort.Strings(got)
		sort.Strings(tc.expected)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Fatalf("expected %v, but got %v", tc.expected, got)
		}
	}
}
