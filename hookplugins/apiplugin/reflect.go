package apiplugin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alibaba/pouch/apis/types"
	hp "github.com/alibaba/pouch/hookplugins"
)

type convertToHostConfigFieldFunc func(string) (interface{}, error)

func convertToString(data string) (interface{}, error) {
	return data, nil
}

func convertToArrayString(data string) (interface{}, error) {
	return strings.Split(data, " "), nil
}

func convertStringToInt(data string) (interface{}, error) {
	v, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return int(0), fmt.Errorf("failed to convert %s to int: %v", data, err)
	}

	return int(v), nil
}

func convertStringToInt64(data string) (interface{}, error) {
	v, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return int64(0), fmt.Errorf("failed to convert %s to int64: %v", data, err)
	}

	return v, nil
}

func convertStringToIntPtr(data string) (interface{}, error) {
	v, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s to int: %v", data, err)
	}

	value := int(v)
	return &value, nil
}

func convertStringToSliceThrottle(data string) (interface{}, error) {
	if data == "" {
		return nil, nil
	}

	tds := []*types.ThrottleDevice{}

	tdStrings := strings.Split(data, " ")
	for _, tdStr := range tdStrings {
		kvs := strings.SplitN(tdStr, ":", 2)
		if len(kvs) != 2 {
			return nil, fmt.Errorf("failed to convert %s to sliceThrottle", data)
		}

		rate, err := strconv.ParseUint(kvs[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert %s to sliceThrottle", data)
		}

		td := &types.ThrottleDevice{
			Path: kvs[0],
			Rate: rate,
		}

		tds = append(tds, td)
	}

	return tds, nil
}

// convertInfo describes
type convertInfo struct {
	fieldName   string
	convertFunc convertToHostConfigFieldFunc
}

func newConvertInfo(fieldName string, f convertToHostConfigFieldFunc) *convertInfo {
	return &convertInfo{
		fieldName:   fieldName,
		convertFunc: f,
	}
}

// resourceWrapReflectMap return a map which maps annotation key to reflect info, which store field name and convertFunc. convertFunc will convert string to hostConfig field value
func resourceWrapReflectMap() map[string]*convertInfo {
	return map[string]*convertInfo{
		hp.SpecCpusetTrickCpus:        newConvertInfo("CpusetTrickCpus", convertToString),
		hp.SpecCpusetTrickTasks:       newConvertInfo("CpusetTrickTasks", convertToString),
		hp.SpecCpusetTrickExemptTasks: newConvertInfo("CpusetTrickExemptTasks", convertToString),
		hp.SpecCPUBvtWarpNs:           newConvertInfo("CPUBvtWarpNs", convertStringToInt64),
		hp.SpecCpuacctSchedLatSwitch:  newConvertInfo("ScheLatSwitch", convertStringToInt64),

		hp.SpecMemoryWmarkRatio:     newConvertInfo("MemoryWmarkRatio", convertStringToInt),
		hp.SpecMemoryExtraInBytes:   newConvertInfo("MemoryExtra", convertStringToInt64),
		hp.SpecMemoryForceEmptyCtl:  newConvertInfo("MemoryForceEmptyCtl", convertStringToInt),
		hp.SpecMemoryPriority:       newConvertInfo("MemoryPriority", convertStringToIntPtr),
		hp.SpecMemoryUsePriorityOOM: newConvertInfo("MemoryUsePriorityOOM", convertStringToIntPtr),
		hp.SpecMemoryOOMKillAll:     newConvertInfo("MemoryKillAll", convertStringToIntPtr),
		hp.SpecMemoryDroppable:      newConvertInfo("MemoryDroppable", convertStringToInt),

		hp.SpecIntelRdtL3Cbm: newConvertInfo("IntelRdtL3Cbm", convertToString),
		hp.SpecIntelRdtGroup: newConvertInfo("IntelRdtGroup", convertToString),
		hp.SpecIntelRdtMba:   newConvertInfo("IntelRdtMba", convertToString),

		hp.SpecBlkioFileLevelSwitch:      newConvertInfo("BlkFileLevelSwitch", convertStringToInt),
		hp.SpecBlkioBufferWriteBps:       newConvertInfo("BlkBufferWriteBps", convertStringToInt),
		hp.SpecBlkioMetaWriteTps:         newConvertInfo("BlkMetaWriteTps", convertStringToInt),
		hp.SpecBlkioFileThrottlePath:     newConvertInfo("BlkFileThrottlePath", convertToArrayString),
		hp.SpecBlkioBufferWriteSwitch:    newConvertInfo("BlkBufferWriteSwitch", convertStringToInt),
		hp.SpecBlkioDeviceBufferWriteBps: newConvertInfo("BlkDeviceBufferWriteBps", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceIdleTime:       newConvertInfo("BlkDeviceIdleTime", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceLatencyTarget:  newConvertInfo("BlkDeviceLatencyTarget", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceReadLowBps:     newConvertInfo("BlkioDeviceReadLowBps", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceReadLowIOps:    newConvertInfo("BlkioDeviceReadLowIOps", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceWriteLowBps:    newConvertInfo("BlkioDeviceWriteLowBps", convertStringToSliceThrottle),
		hp.SpecBlkioDeviceWriteLowIOps:   newConvertInfo("BlkioDeviceWriteLowIOps", convertStringToSliceThrottle),

		hp.SpecNetCgroupRate: newConvertInfo("NetCgroupRate", convertToString),
		hp.SpecNetCgroupCeil: newConvertInfo("NetCgroupCeil", convertToString),
	}
}
