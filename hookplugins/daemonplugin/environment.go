package daemonplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/alibaba/pouch/pkg/kernel"
)

const hotfixesPath = "/opt/ali-iaas/pouch/hotfixesList.json"

type consistencyConfig struct {
	Kernel string `json:"kernel,omitempty"`
	Value  string `json:"value,omitempty"`
}

// check the consistency of environment
func checkEnvConsistency(pouchRootDir string) error {
	// 1. check the kernel version
	kernelVersion := "unknown"
	version, err := kernel.GetKernelVersion()
	if err != nil {
		return fmt.Errorf("Could not get kernel version: %v", err)
	}
	kernelVersion = version.String()

	// 2. check the hotfixes
	if err := checkHotfixes(kernelVersion); err != nil {
		return err
	}

	// 3. check dirquota
	if err := checkDirquota(kernelVersion, pouchRootDir); err != nil {
		return err
	}

	return nil
}

func checkHotfixes(kernelVersion string) error {
	hotfixes := make(map[string][]string)

	// parse hotfixesList.json
	data, err := ioutil.ReadFile(hotfixesPath)
	if err != nil {
		return fmt.Errorf("read hotfixesList.json error, %v", err)
	}
	v := []consistencyConfig{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("parse the hotfixesList.json error, %v", err)
	}
	for _, hotfix := range v {
		splits := strings.Split(hotfix.Value, ",")
		for _, val := range splits {
			hotfixes[hotfix.Kernel] = append(hotfixes[hotfix.Kernel], val)
		}
	}

	// check the hotfixes corresponding to kernel version
	if hotfixes[kernelVersion] == nil {
		return fmt.Errorf("unknown kernel version")
	}
	hotfixesPath := "/lib/modules/" + kernelVersion + "/extra/pouchhotfixes/"
	for _, val := range hotfixes[kernelVersion] {
		tempPath := hotfixesPath + val
		_, err := os.Stat(tempPath)
		if err != nil {
			return fmt.Errorf("hotfix %v does not exist", tempPath)
		}
	}

	return nil
}

func checkDirquota(kernelVersion, pouchRootDir string) error {

	// parse kernel version
	str := strings.Split(kernelVersion, ".")
	version := str[0] + str[1]

	// get graph device
	device, err := getGraphDevice(pouchRootDir)
	if err != nil {
		return err
	}

	// according to the kernel version, we should check separately.
	err = isQuotaOn(version, device)
	if err != nil {
		return err
	}

	return nil
}

func execCmd(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("run command error, %v", err)
	}

	return out.String(), nil
}

func getGraphDevice(pouchRootDir string) (string, error) {
	var device string
	deviceData, err := execCmd("df", pouchRootDir)
	if err != nil {
		return "", fmt.Errorf("get device error, %v", err)
	}
	for _, line := range strings.Split(deviceData, "\n") {
		if strings.Contains(line, "/dev") {
			temp := strings.Split(line, " ")
			device = temp[0]
			break
		}
	}
	if device == "" {
		return "", fmt.Errorf("can not get graph device")
	}

	return device, nil
}

func isQuotaOn(version, device string) error {
	if version == "49" {
		var features string
		featuresData, err := execCmd("dumpe2fs", "-h", device)
		if err != nil || features == "" {
			return fmt.Errorf("can not get filesystem features")
		}
		for _, line := range strings.Split(featuresData, "\n") {
			if strings.Contains(line, "Filesystem features") {
				features = line
				break
			}
		}
		if !strings.Contains(features, "project") {
			return fmt.Errorf("no enable project fs feature")
		}
		if !strings.Contains(features, "quota") {
			return fmt.Errorf("no enable quota fs feature")
		}
	} else {
		grpquota := false
		grpquotaData, err := execCmd("mount", "-l")
		if err != nil {
			return fmt.Errorf("get mount info error, %v", err)
		}
		for _, line := range strings.Split(grpquotaData, "\n") {
			if strings.Contains(line, device) && strings.Contains(line, "grpquota") {
				grpquota = true
				break
			}
		}
		if !grpquota {
			return fmt.Errorf("no enable grpquota")
		}
	}
	return nil
}
