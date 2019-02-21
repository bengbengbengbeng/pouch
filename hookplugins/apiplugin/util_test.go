package apiplugin

import (
	"testing"

	"github.com/alibaba/pouch/apis/types"
)

func TestSliceThrottleDeviceString(t *testing.T) {
	var (
		tdsSlice1 []*types.ThrottleDevice
		tdsSlice2 []*types.ThrottleDevice
	)

	tds1 := &types.ThrottleDevice{
		"abc",
		1,
	}
	tds2 := &types.ThrottleDevice{
		"def",
		2,
	}
	tdsSlice1 = append(tdsSlice1, tds1)
	tdsSlice2 = append(tdsSlice2, tds1, tds2)

	cases := []struct {
		tds      []*types.ThrottleDevice
		expected string
	}{
		{
			tds:      nil,
			expected: "",
		}, {
			tds:      tdsSlice1,
			expected: "abc:1",
		},
		{
			tds:      tdsSlice2,
			expected: "abc:1 def:2",
		},
	}

	for _, tc := range cases {
		got := sliceThrottleDeviceString(tc.tds)
		if got != tc.expected {
			t.Fatalf("expected %v, but got %v", tc.expected, got)
		}
	}
}
