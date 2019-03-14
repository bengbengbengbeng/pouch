package main

import (
	"os/exec"
	"strings"
	"time"

	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// PouchPanguImageSuite is the test suite for Pangu image.
type PouchPanguImageSuite struct{}

func init() {
	check.Suite(&PouchPanguImageSuite{})
}

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchPanguImageSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)

	environment.PruneAllContainers(apiClient)

	startTDCService(c)
}

// TearDownTest does cleanup work in the end of each test.
func (suite *PouchPanguImageSuite) TearDownTest(c *check.C) {
	environment.PruneAllContainers(apiClient)
}

// TestRestartDaemon in TODO list.
func (suite *PouchPanguImageSuite) TestRestartDaemon(c *check.C) {
	helpwantedForMissingCase(c, "running container with pangu image and restart pouchd")
}

// TestPanguImageBasic will check the basical functionality of pangu image.
//
// NOTE: The remote disk must not be created during pulling an pangu image, and
// the remote disk can only be created for container snapshotter and only
// unmounted after children overlay snapshotter has been removed because the
// container can be restart again.
//
// The id of remote disk is random so that it is hard to check mount info
// precisely. Use the number of remote disk is one way so that the case cannot
// run parallelly.
func (suite *PouchPanguImageSuite) TestPanguImageBasic(c *check.C) {
	cname := "TestPanguImageBasic"
	panguImage := "reg.docker.alibaba-inc.com/pangu/test_ali_os:7u2"

	command.PouchRun("pull", panguImage).Assert(c, icmd.Success)

	// download image should not create any remote disk in local
	checkPanguRDBMountPointNumber(c, 0, false)

	// run container in detach mode
	res := command.PouchRun("run", "-d",
		"--name", cname,
		panguImage, "sleep", "10000")
	defer DelContainerForceMultyTime(c, cname)
	res.Assert(c, icmd.Success)

	// start container action will need only one remote disk
	checkPanguRDBMountPointNumber(c, 1, false)

	// create /tmp/1 file in container
	command.PouchRun("exec", cname, "sh", "-c", "echo hi > /tmp/1").Assert(c, icmd.Success)

	// stop container but we still need the remote disk in local
	command.PouchRun("stop", cname).Assert(c, icmd.Success)
	checkPanguRDBMountPointNumber(c, 1, false)

	// restart the stopped container and check the /tmp/1 file
	command.PouchRun("start", cname).Assert(c, icmd.Success)
	res = command.PouchRun("exec", cname, "sh", "-c", "cat /tmp/1")
	res.Assert(c, icmd.Success)
	if got := strings.TrimSpace(res.Combined()); got != "hi" {
		c.Fatalf("expected to get hi, but got %v", got)
	}

	// after remove container, the remote disk should be remove
	command.PouchRun("stop", cname).Assert(c, icmd.Success)
	checkPanguRDBMountPointNumber(c, 1, false)
	command.PouchRun("rm", cname).Assert(c, icmd.Success)

	// NOTE: the remote disk will be removed by contained gc.
	// the check will be retry in about 10 sec.
	checkPanguRDBMountPointNumber(c, 0, true)
}

// checkPanguRDBMountPoint will get mount info to check the rbd mount point
func checkPanguRDBMountPointNumber(t testingTB, n int, retry bool) {
	var (
		common = "io.containerd.snapshotter.v1.rbd/rdisks/pangu"
		total  = 10
		got    = 0
	)

	for i := 0; i < total; i++ {
		cmd := exec.Command("df")
		data, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("failed to get mount info: %v", err)
		}

		got = strings.Count(string(data), common)
		if got == n {
			return
		}

		if retry {
			time.Sleep(1 * time.Second)
		}
	}

	if got != n {
		t.Fatalf("expected rbd mount number(%v), but got number(%v)", n, got)
	}
}

func startTDCService(t testingTB) {
	cmd := exec.Command("service", "tdc", "restart")
	if _, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to start tdc service: %v", err)
	}

	// NOTE: there is something wrong with tdc service.
	// we cann't use systemctl show -p ActiveState tdc to get real status
	// so that we can only do it by strings matching.
	cmd = exec.Command("service", "tdc", "status")
	data, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to get tdc service status: %v", err)
	}

	// there are two daemons here: watchdog and tdc
	expected := 2
	if got := strings.Count(string(data), "[OK]"); got != expected {
		t.Fatalf("expected to active all the services, but got %v\n", string(data))
	}
}
