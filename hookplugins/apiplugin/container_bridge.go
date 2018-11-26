package apiplugin

import (
	"context"
	"net/http"
	"runtime"
	"strings"

	"github.com/alibaba/pouch/apis/server"
	serverTypes "github.com/alibaba/pouch/apis/server/types"
	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/pkg/kernel"
	"github.com/alibaba/pouch/pkg/utils"
	"github.com/alibaba/pouch/version"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

func getVersionHandler(_ serverTypes.Handler) serverTypes.Handler {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) (err error) {
		kernelVersion := "<unknown>"
		if kv, err := kernel.GetKernelVersion(); err != nil {
			logrus.Warnf("Could not get kernel version: %v", err)
		} else {
			kernelVersion = kv.String()
		}

		v := types.SystemVersion{
			APIVersion:    version.APIVersion,
			Arch:          runtime.GOARCH,
			BuildTime:     version.BuildTime,
			GitCommit:     version.GitCommit,
			GoVersion:     runtime.Version(),
			KernelVersion: kernelVersion,
			Os:            runtime.GOOS,
			Version:       version.Version,
		}

		if utils.IsStale(ctx, req) {
			v.Version = "1.12.6"
		}

		return server.EncodeResponse(w, http.StatusOK, v)
	}
}

func containerRequestWrapper(h serverTypes.Handler) serverTypes.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		if utils.IsStale(ctx, req) {
			//trim the heading slash in parameter name
			nameInPath := mux.Vars(req)["name"]
			if strings.HasPrefix(nameInPath, "/") {
				nameInPath = strings.TrimPrefix(nameInPath, "/")
				mux.Vars(req)["name"] = nameInPath
			}

			nameInForm := req.FormValue("name")
			if strings.HasPrefix(nameInForm, "/") {
				nameInForm = strings.TrimPrefix(nameInForm, "/")
				req.Form.Set("name", nameInForm)
			}
		}

		return h(ctx, rw, req)
	}
}
