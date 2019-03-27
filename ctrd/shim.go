package ctrd

import (
	"context"
	"path/filepath"

	"github.com/alibaba/pouch/ctrd/supervisord"

	"github.com/containerd/containerd/runtime/v1/shim"
	shimclient "github.com/containerd/containerd/runtime/v1/shim/client"
	"github.com/sirupsen/logrus"
)

// NewShimClient is going to connect to a container process'shim
func NewShimClient(ctx context.Context, cID, namespace string) (*shimclient.Client, error) {
	cfg, err := generateShimConfig(cID, namespace)
	if err != nil {
		return nil, err
	}

	opt := shimclient.WithConnect(shimAddress(namespace, cID), func() {
		logrus.Infof("ctrd: containerd-shim exited: %s", cID)
	})

	return shimclient.New(ctx, *cfg, opt)
}

func shimAddress(ns, id string) string {
	return filepath.Join(string(filepath.Separator), "containerd-shim", ns, id, "shim.sock")
}

func generateShimConfig(cID, namespace string) (*shim.Config, error) {
	ctrdCfg, err := supervisord.DaemonConfig()
	if err != nil {
		return nil, err
	}

	cfg := &shim.Config{
		Namespace:   namespace,
		RuntimeRoot: runtimeRoot,
		Path:        filepath.Join(ctrdCfg.Root, "io.containerd.runtime.v1.linux", cID),
		WorkDir:     filepath.Join(ctrdCfg.State, "io.containerd.runtime.v1.linux", cID),
		// SystemdCgroup: false,
		// Criu:          "",
	}

	return cfg, nil

}
