package containerplugin

import (
	"net"
	"strings"
	"time"

	networktypes "github.com/alibaba/pouch/network/types"

	"github.com/sirupsen/logrus"
)

// PreStart returns an array of priority and args which will pass to runc, the every priority
// used to sort the pre start array that pass to runc, network plugin hook always has priority value 0.
// Prestart copy files to container rootfs
func (c *contPlugin) PreStart(interface{}) ([]int, [][]string, error) {
	logrus.Infof("pre start method called")

	// invoke script at /opt/ali-iaas/pouch/bin/start_hook.sh
	// copy file into the ns, put entrypoint in container. like the function of pouch_container_create.sh in old version
	return []int{-100}, [][]string{{"/opt/ali-iaas/pouch/bin/prestart_hook"}}, nil
}

// PreCreateEndpoint accepts the container id and env of this container, to update the config of container's endpoint.
// 1. pass Overlay parameters to network plugin like alinet
// 2. generate mac address from ip address
// 3. generate priority for the network interface
func (c *contPlugin) PreCreateEndpoint(cid string, env []string, endpoint *networktypes.Endpoint) error {
	genericParam := make(map[string]interface{})
	if getEnv(env, "OverlayNetwork") == "true" {
		genericParam["OverlayNetwork"] = "true"
		genericParam["OverlayTunnelId"] = getEnv(env, "OverlayTunnelId")
		genericParam["OverlayGwIp"] = getEnv(env, "OverlayGwIp")
	}

	if getEnv(env, "VpcECS") == "true" {
		genericParam["VpcECS"] = "true"
	}

	for _, oneEnv := range env {
		arr := strings.SplitN(oneEnv, "=", 2)
		if len(arr) == 2 && strings.HasPrefix(arr[0], "alinet_") {
			genericParam[arr[0]] = arr[1]
		}
	}

	if ip := getEnv(env, "RequestedIP"); ip != "" {
		if strings.Contains(ip, ",") {
			ip = strings.Split(ip, ",")[0]
		}
		genericParam[MacAddress] = GenerateMACFromIP(net.ParseIP(ip))
	}

	endpoint.Priority = int(finalPoint.Unix() - time.Now().Unix())
	endpoint.DisableResolver = true
	endpoint.GenericParams = genericParam

	return nil
}
