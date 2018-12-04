package daemonplugin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/alibaba/pouch/apis/metrics"
	"github.com/alibaba/pouch/hookplugins"

	"github.com/sirupsen/logrus"
)

type daemonPlugin struct{}

func init() {
	hookplugins.RegisterDaemonPlugin(&daemonPlugin{})
}

// PreStartHook is invoked by pouch daemon before real start, in this hook user could start http proxy or other
// standalone process plugins
// copy config from /etc/pouch/config.json to /etc/sysconfig/pouch
// and start plugin processes than daemon depended
func (d *daemonPlugin) PreStartHook() error {
	logrus.Infof("pre-start hook in daemon is called")
	configMap := make(map[string]interface{}, 8)
	if _, ex := os.Stat("/etc/pouch/config.json"); ex == nil {
		f, err := os.OpenFile("/etc/pouch/config.json", os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		err = json.NewDecoder(f).Decode(&configMap)
		f.Close()
		if err != nil {
			return err
		}
		if s, ok := configMap["home-dir"].(string); ok {
			homeDir = s
		}
	}
	setupEnv()
	b, e := exec.Command("/opt/ali-iaas/pouch/bin/daemon_prestart.sh", homeDir).CombinedOutput()
	if e != nil {
		return fmt.Errorf("daemon prestart execute error. %s %v", string(b), e)
	}
	logrus.Infof("daemon_prestart output %s", string(b))

	// check the consistency of environment
	status := 0
	if err := checkEnvConsistency(homeDir); err != nil {
		logrus.Errorf("the environment is not consistent, %v", err)
		metrics.EnvStatus.WithLabelValues("consistency").Set(1)
		status = 1
	} else {
		metrics.EnvStatus.WithLabelValues("consistency").Set(0)
	}
	logrus.Infof("check the consistency of environment, %v", status)
	go activePlugins()
	return nil
}

// PreStopHook stops plugin processes than start ed by PreStartHook.
func (d *daemonPlugin) PreStopHook() error {
	logrus.Infof("pre-stop hook in daemon is called")
	b, e := exec.Command("/opt/ali-iaas/pouch/bin/daemon_prestop.sh").CombinedOutput()
	if e != nil {
		return e
	}
	logrus.Infof("daemon_prestop output %s", string(b))
	return nil
}
