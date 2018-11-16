package daemonplugin

import (
	"os"
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

func TestCheckHotfixesSuccess(t *testing.T) {
	kernelVersion := "3.10.0-327.ali2000.alios7.x86_64"
	path := "/lib/modules/" + kernelVersion
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

	assert.NilError(t, err)
}

func TestCheckHotfixesFailHotfixesNotExist(t *testing.T) {
	kernelVersion := "3.10.0-327.ali2000.alios7.x86_64"
	err := checkHotfixes(kernelVersion)

	assert.Error(t, err, "does not exist")
}

func TestCheckHotfixesFailUnknownKernelVersion(t *testing.T) {
	kernelVersion := "1.3.10.0-327.ali2000.alios7.x86_64"
	err := checkHotfixes(kernelVersion)

	assert.Error(t, err, "unknown kernel version")
}

func TestCheckDirQuota310NotExist(t *testing.T) {
	kernelVersion := "3.10.0-327.ali2000.alios7.x86_64"
	pouchRootDir := "/home/t4/pouch"
	err := checkDirquota(kernelVersion, pouchRootDir)

	assert.Error(t, err, "no enable grpquota")
}
