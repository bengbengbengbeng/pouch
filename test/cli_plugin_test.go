package main

import (
	"strings"

	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// PouchPluginSuite is the test suite for ps CLI.
type PouchPluginSuite struct{}

func init() {
	check.Suite(&PouchPluginSuite{})
}

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchPluginSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)

	environment.PruneAllContainers(apiClient)
	PullImage(c, busyboxImage)
}

// TearDownTest does cleanup work in the end of each test.
func (suite *PouchPluginSuite) TearDownTest(c *check.C) {
}

func (suite *PouchRunSuite) TestRunQuotaId(c *check.C) {
	if !environment.IsDiskQuota() {
		c.Skip("Host does not support disk quota")
	}
	name := "TestRunQuotaId"
	ID := "16777216"

	res := command.PouchRun("run", "-d", "--name", name, "--label", "DiskQuota=10G", "--label", "QuotaId="+ID, busyboxImage, "top")
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	output := command.PouchRun("inspect", "-f", "{{.Config.Labels.QuotaId}}", name).Stdout()
	c.Assert(strings.TrimSpace(output), check.Equals, ID)

}

func (suite *PouchRunSuite) TestRunAutoQuotaId(c *check.C) {
	if !environment.IsDiskQuota() {
		c.Skip("Host does not support disk quota")
	}
	name := "TestRunAutoQuotaId"
	AutoQuotaIDValue := "true"

	res := command.PouchRun("run", "-d", "--name", name, "--label", "DiskQuota=10G", "--label", "AutoQuotaId="+AutoQuotaIDValue, busyboxImage, "top")
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	output := command.PouchRun("inspect", "-f", "{{.Config.Labels.AutoQuotaId}}", name).Stdout()
	c.Assert(strings.TrimSpace(output), check.Equals, AutoQuotaIDValue)
}

// TestRunDiskQuotaForAllDirsWithoutQuotaId: quota-id(<0) disk-quota:.*=10GB
func (suite *PouchRunSuite) TestRunDiskQuotaForAllDirsWithoutQuotaId(c *check.C) {
	if !environment.IsDiskQuota() {
		c.Skip("Host does not support disk quota")
	}
	name := "TestRunDiskQuotaForAllDirsWithoutQuotaId"

	volumeName1 := "volume1"
	volumeName2 := "volume2"

	command.PouchRun("volume", "create", "-n", volumeName1, "-d", "local").Assert(c, icmd.Success)
	defer command.PouchRun("volume", "rm", volumeName1)

	command.PouchRun("volume", "create", "-n", volumeName2, "-d", "local", "-o", "opt.size=30M").Assert(c, icmd.Success)
	defer command.PouchRun("volume", "rm", volumeName2)

	command.PouchRun("run", "-d", `--disk-quota=".*=1g"`, "--quota-id=-1", "-v", volumeName1+":/mnt1", "-v", volumeName2+":/mnt2", "--name", name, busyboxImage, "top").Assert(c, icmd.Success)
	defer command.PouchRun("rm", "-f", name)

	expct := icmd.Expected{
		ExitCode: 1,
		Err:      "Disk quota exceeded",
	}
	cmd := "dd if=/dev/zero of=/mnt2/test1 bs=1M count=20"
	command.PouchRun("exec", name, "sh", "-c", cmd).Assert(c, icmd.Success)

	cmd = "dd if=/dev/zero of=/mnt2/test2 bs=1M count=20"
	err := command.PouchRun("exec", name, "sh", "-c", cmd).Compare(expct)
	c.Assert(err, check.IsNil)

	expectedstring := "1048576"
	cmd = "df | grep mnt1"
	out := command.PouchRun("exec", name, "sh", "-c", cmd).Stdout()
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}
}
