package apiplugin

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alibaba/pouch/apis/types"
	hp "github.com/alibaba/pouch/hookplugins"
)

func TestSliceThrottleDeviceString(t *testing.T) {
	var (
		tdsSlice1 []*types.ThrottleDevice
		tdsSlice2 []*types.ThrottleDevice
	)

	tds1 := &types.ThrottleDevice{
		Path: "abc",
		Rate: 1,
	}
	tds2 := &types.ThrottleDevice{
		Path: "def",
		Rate: 2,
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

func TestReflect(t *testing.T) {
	resource := &ResourcesWrapper{}
	cputrickCpus := "1,2,3,4"
	memoryWmarkRatio := 10
	memoryExtra := int64(1000)

	memoryPriorityValue := 100
	memoryPriority := &memoryPriorityValue
	blkDeviceBufferWriteBps := []*types.ThrottleDevice{
		{
			Path: "/test1",
			Rate: uint64(100),
		},
	}

	r := reflect.ValueOf(resource).Elem()
	r.FieldByName("CpusetTrickCpus").Set(reflect.ValueOf(cputrickCpus))
	r.FieldByName("MemoryWmarkRatio").Set(reflect.ValueOf(memoryWmarkRatio))
	r.FieldByName("MemoryExtra").Set(reflect.ValueOf(memoryExtra))
	r.FieldByName("MemoryPriority").Set(reflect.ValueOf(memoryPriority))
	r.FieldByName("BlkDeviceBufferWriteBps").Set(reflect.ValueOf(blkDeviceBufferWriteBps))

	if resource.CpusetTrickCpus != cputrickCpus {
		t.Fatalf("CpusetTrickCpus expected %v, but %v", cputrickCpus, resource.CpusetTrickCpus)
	}

	if resource.MemoryWmarkRatio != memoryWmarkRatio {
		t.Fatalf("MemoryWmarkRatio expected %v, but %v", memoryWmarkRatio, resource.MemoryWmarkRatio)
	}

	if resource.MemoryExtra != memoryExtra {
		t.Fatalf("MemoryExtra expected %v, but %v", memoryWmarkRatio, resource.MemoryExtra)
	}

	if *resource.MemoryPriority != *memoryPriority {
		t.Fatalf("MemoryPriority expected %v, but %v", *memoryPriority, *resource.MemoryPriority)
	}

	if resource.BlkDeviceBufferWriteBps[0].Path != blkDeviceBufferWriteBps[0].Path &&
		resource.BlkDeviceBufferWriteBps[0].Rate != blkDeviceBufferWriteBps[0].Rate {
		t.Fatalf("BlkDeviceBufferWriteBps expected %v, but %v", blkDeviceBufferWriteBps, resource.BlkDeviceBufferWriteBps)
	}
}

func TestConvertAnnotationToHostConfig(t *testing.T) {
	annotations := map[string]string{
		hp.SpecCpusetTrickCpus:       "1,2",
		hp.SpecMemoryWmarkRatio:      "1",
		hp.SpecMemoryExtraInBytes:    "1024",
		hp.SpecMemoryPriority:        "2",
		hp.SpecBlkioFileThrottlePath: "/path1 /path2",
		hp.SpecBlkioDeviceReadLowBps: "/path1:1024 /path2:2048",
	}

	resource := &ResourcesWrapper{}
	err := convertAnnotationToDockerHostConfig(annotations, resource)
	if err != nil {
		t.Fatalf("expected not nil, but %v", err)
	}

	if resource.CpusetTrickCpus != "1,2" {
		t.Fatalf("CpusetTrickCpus expected to be 1,2, but get %v", resource.CpusetTrickCpus)
	}

	if resource.MemoryWmarkRatio != 1 {
		t.Fatalf("MemoryWmarkRatio expected to be 1, but get %v", resource.MemoryWmarkRatio)
	}

	if resource.MemoryExtra != int64(1024) {
		t.Fatalf("MemoryExtra expected to be 1024, but get %v", resource.MemoryExtra)
	}

	if *resource.MemoryPriority != 2 {
		t.Fatalf("MemoryPriority expected to be 2, but get %v", *resource.MemoryPriority)
	}

	blkReadLowBps := resource.BlkioDeviceReadLowBps

	if len(blkReadLowBps) != 2 || fmt.Sprintf("%s:%d", blkReadLowBps[0].Path, blkReadLowBps[0].Rate) != "/path1:1024" || fmt.Sprintf("%s:%d", blkReadLowBps[1].Path, blkReadLowBps[1].Rate) != "/path2:2048" {
		t.Fatalf("BlkioDeviceReadLowBps expected to be /path1:1024 /path2:2048, but get %v", resource.BlkioDeviceReadLowBps)
	}
}
