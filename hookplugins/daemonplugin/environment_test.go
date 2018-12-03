package daemonplugin

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckHotfixesSuccess(t *testing.T) {
	kernelVersion := "3.10.0-327.ali2000.alios7.x86_64"
	path := "/lib/modules/" + kernelVersion + "/extra/pouchhotfixes/"
	cgroupMount := path + "/cgroup_mount"
	if _, err := os.Stat(cgroupMount); err != nil {
		err := os.MkdirAll(cgroupMount, os.ModePerm)
		if err != nil {
			t.Errorf("mkdir %v error, %v", cgroupMount, err)
		}
	}
	hotfixInfo := path + "/hotfix_info"
	if _, err := os.Stat(hotfixInfo); err != nil {
		err := os.MkdirAll(hotfixInfo, os.ModePerm)
		if err != nil {
			t.Errorf("mkdir %v error, %v", hotfixInfo, err)
		}
	}
	vhost := path + "/vhost_net"
	if _, err := os.Stat(vhost); err != nil {
		err := os.MkdirAll(vhost, os.ModePerm)
		if err != nil {
			t.Errorf("mkdir %v error, %v", vhost, err)
		}
	}
	err := checkHotfixes(kernelVersion)

	if err := os.RemoveAll(path); err != nil {
		t.Errorf("remove %v error, %v", path, err)
	}

	assert.NoError(t, err)
}

func TestCheckHotfixesFailHotfixesNotExist(t *testing.T) {
	kernelVersion := "3.10.0-327.ali2000.alios7.x86_64"
	err := checkHotfixes(kernelVersion)

	assert.Equal(t, "hotfix /lib/modules/3.10.0-327.ali2000.alios7.x86_64/extra/pouchhotfixes/cgroup_mount does not exist", err.Error())
}

func TestCheckHotfixesFailUnknownKernelVersion(t *testing.T) {
	kernelVersion := "1.3.10.0-327.ali2000.alios7.x86_64"
	err := checkHotfixes(kernelVersion)

	assert.Equal(t, "unknown kernel version", err.Error())
}
