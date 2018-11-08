package utils

import (
	"context"
	"net/http"
	"strings"
)

// IsStale checks tls name and return if this context is owned by sigma or use docker client
func IsStale(ctx context.Context, req *http.Request) bool {
	isRemoteSigma := GetTLSIssuer(ctx) == "ali" && strings.Contains(GetTLSCommonName(ctx), "sigma")
	if req == nil {
		return isRemoteSigma
	}
	return isRemoteSigma || strings.Contains(req.Header.Get("User-Agent"), "Docker-Client")
}
