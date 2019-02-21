package apiplugin

import (
	"fmt"
	"strings"

	"github.com/alibaba/pouch/apis/types"
)

// isContainerRequest check if the request is a container related request
func isContainerRequest(path string) bool {
	return strings.HasPrefix(path, "/containers/")
}

// sliceThrottleDeviceString formats struct ThrottleDevice to string
func sliceThrottleDeviceString(tds []*types.ThrottleDevice) string {
	if tds == nil {
		return ""
	}

	res := make([]string, 0, len(tds))
	for _, t := range tds {
		if t == nil {
			continue
		}
		res = append(res, fmt.Sprintf("%s:%d", t.Path, t.Rate))
	}
	return strings.Join(res, " ")
}
