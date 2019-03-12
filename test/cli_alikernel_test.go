package main

import (
	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/util"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// PouchAliKernelSuite is the test suite for AliOS specified features.
type PouchAliKernelSuite struct{}

func init() {
	check.Suite(&PouchAliKernelSuite{})
}

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchAliKernelSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
	SkipIfFalse(c, environment.IsAliKernel)

	environment.PruneAllContainers(apiClient)

	PullImage(c, busyboxImage)
}

// TearDownTest does cleanup work in the end of each test.
func (suite *PouchAliKernelSuite) TearDownTest(c *check.C) {
}

// TestAliKernelDiskQuotaWorks tests disk quota works on AliKernel.
func (suite *PouchAliKernelSuite) TestAliKernelDiskQuotaWorks(c *check.C) {
	funcname := "TestAliKernelDiskQuotaWorks"

	command.PouchRun("volume", "create", "--name", funcname, "-d", "local", "-o", "opt.size=1g").Assert(c, icmd.Success)
	defer command.PouchRun("volume", "rm", funcname)

	command.PouchRun("run", "-d", "-v", funcname+":/mnt", "--name", funcname, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, funcname)

	expct := icmd.Expected{
		ExitCode: 0,
		Out:      "1.0G",
	}
	err := command.PouchRun("exec", funcname, "df", "-h").Compare(expct)
	c.Assert(err, check.IsNil)

	// generate a file larger than 1G should fail.
	expct = icmd.Expected{
		ExitCode: 1,
		Err:      "Disk quota exceeded",
	}
	cmd := "dd if=/dev/zero of=/mnt/test bs=1024k count=1500"
	err = command.PouchRun("exec", funcname, "sh", "-c", cmd).Compare(expct)
	c.Assert(err, check.IsNil)
}

// TestAliKernelDiskQuotaMultiWorks tests multi volume with different disk quota works on AliKernel.
func (suite *PouchAliKernelSuite) TestAliKernelDiskQuotaMultiWorks(c *check.C) {
	funcname := "TestAliKernelDiskQuotaMultiWorks"

	name1 := funcname + "1"
	name2 := funcname + "2"

	command.PouchRun("volume", "create", "--name", name1, "-d", "local", "-o", "opt.size=2.2g").Assert(c, icmd.Success)
	defer command.PouchRun("volume", "rm", name1)

	command.PouchRun("volume", "create", "--name", name2, "-d", "local", "-o", "opt.size=3.2g").Assert(c, icmd.Success)
	defer command.PouchRun("volume", "rm", name2)

	command.PouchRun("run", "-d", "-v", name1+":/mnt/test1", "-v", name2+":/mnt/test2", "--name", funcname, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, funcname)

	{
		expct := icmd.Expected{
			ExitCode: 0,
			Out:      "2.2G",
		}
		cmd := "df -h |grep test1"
		err := command.PouchRun("exec", funcname, "sh", "-c", cmd).Compare(expct)
		c.Assert(err, check.IsNil)
	}
	{
		expct := icmd.Expected{
			ExitCode: 0,
			Out:      "3.2G",
		}
		cmd := "df -h |grep test2"
		err := command.PouchRun("exec", funcname, "sh", "-c", cmd).Compare(expct)
		c.Assert(err, check.IsNil)
	}
}

// TestAliKernelCgroupNS test container cgroup view should be / on both 3.10 and 4.9 alios
func (suite *PouchAliKernelSuite) TestAliKernelCgroupNS(c *check.C) {
	name1 := "TestAliKernelCgroupNS-run"
	name2 := "TestAliKernelCgroupNS-exec"
	defer DelContainerForceMultyTime(c, name1)
	defer DelContainerForceMultyTime(c, name2)

	res := command.PouchRun("run", "--name", name1, busyboxImage, "cat", "/proc/self/cgroup")
	res.Assert(c, icmd.Success)

	cgroupPaths := util.ParseCgroupFile(res.Stdout())
	for _, v := range cgroupPaths {
		if v != "/" {
			c.Fatalf("unexpected cgroup path on alikernel %v", v)
		}
	}

	command.PouchRun("run", "-d", "--name", name2, busyboxImage, "top").Assert(c, icmd.Success)
	res = command.PouchRun("exec", name2, "cat", "/proc/self/cgroup")
	res.Assert(c, icmd.Success)

	cgroupPaths = util.ParseCgroupFile(res.Stdout())
	for _, v := range cgroupPaths {
		if v != "/" {
			c.Fatalf("unexpected cgroup path on alikernel %v", v)
		}
	}
}

// TestUpdateNestCgroup tests update with nest cpu cgroup on alikernel should successful
// which powered by runc
func (suite *PouchAliKernelSuite) TestUpdateNestCpuCgroup(c *check.C) {
	name := "TestUpdateNestCgroup-cpuset"
	defer DelContainerForceMultyTime(c, name)

	command.PouchRun("run", "--name", name, "-d", "--privileged", busyboxImage, "top").Assert(c, icmd.Success)
	command.PouchRun("exec", name, "sh", "-c",
		"mkdir -p /sys/fs/cgroup/cpuset/n1 && echo 1 > /sys/fs/cgroup/cpuset/n1/cpuset.cpus").Assert(c, icmd.Success)

	// update cpu should successful
	command.PouchRun("update", "--cpuset-cpus", "0", name).Assert(c, icmd.Success)
	res := command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/cpuset/cpuset.cpus").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "0\n"), check.IsNil)

	// update memory should successful
	command.PouchRun("update", "-m", "1g", name).Assert(c, icmd.Success)
	res = command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/memory/memory.limit_in_bytes").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "1073741824"), check.IsNil)
}

// TestUpdateNestCgroup tests update with nest memory cgroup on alikernel should successful
// which powered by runc
func (suite *PouchAliKernelSuite) TestUpdateNestMemoryCgroup(c *check.C) {
	name := "TestUpdateNestCgroup-memory"

	command.PouchRun("run", "--name", name, "-d", "--privileged", busyboxImage, "top").Assert(c, icmd.Success)
	command.PouchRun("exec", name, "sh", "-c", "mkdir /sys/fs/cgroup/memory/nest").Assert(c, icmd.Success)

	// update cpu should successful
	command.PouchRun("update", "--cpuset-cpus", "0", name).Assert(c, icmd.Success)
	res := command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/cpuset/cpuset.cpus").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "0\n"), check.IsNil)

	// update memory should successful
	command.PouchRun("update", "-m", "1g", name).Assert(c, icmd.Success)
	res = command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/memory/memory.limit_in_bytes").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "1073741824"), check.IsNil)
}

// TestUpdateNestCgroup tests update with nest device cgroup on alikernel should successful
// which powered by runc
func (suite *PouchAliKernelSuite) TestUpdateNestDeviceCgroup(c *check.C) {
	name := "TestUpdateNestCgroup-device"

	command.PouchRun("run", "--name", name, "-d", "--privileged", busyboxImage, "top").Assert(c, icmd.Success)
	command.PouchRun("exec", name, "sh", "-c", "mkdir /sys/fs/cgroup/devices/nest").Assert(c, icmd.Success)

	// update cpu should successful
	command.PouchRun("update", "--cpuset-cpus", "0", name).Assert(c, icmd.Success)
	res := command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/cpuset/cpuset.cpus").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "0\n"), check.IsNil)

	// update memory should successful
	command.PouchRun("update", "-m", "1g", name).Assert(c, icmd.Success)
	res = command.PouchRun("exec", name, "cat", "/sys/fs/cgroup/memory/memory.limit_in_bytes").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "1073741824"), check.IsNil)
}
