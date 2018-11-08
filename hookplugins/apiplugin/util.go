package apiplugin

import (
	"strings"
)

// isContainerRequest check if the request is a container related request
func isContainerRequest(path string) bool {
	return strings.HasPrefix(path, "/containers/")
}
