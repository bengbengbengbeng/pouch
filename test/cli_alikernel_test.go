package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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
	defer DelContainerForceMultyTime(c, name)

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
	defer DelContainerForceMultyTime(c, name)

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

// TestContainerNS tests container ns is new created, should different from host
func (suite *PouchAliKernelSuite) TestContainerNS(c *check.C) {
	SkipIfFalse(c, func() bool {
		if _, err := exec.LookPath("readlink"); err != nil {
			return false
		}
		return true
	})

	name := "TestContainerNS"
	defer DelContainerForceMultyTime(c, name)
	command.PouchRun("run", "--name", name, "-d", busyboxImage, "top").Assert(c, icmd.Success)

	pid := strings.TrimSpace(command.PouchRun("inspect", "-f", "{{.State.Pid}}", name).Stdout())
	c.Assert(strings.TrimSpace(pid), check.Not(check.Equals), "0")

	getNS := func(pid, ns string) (string, error) {
		rawdata, err := exec.Command("readlink", "/proc/"+pid+"/ns/"+ns).Output()
		out := string(rawdata)
		if err != nil {
			return "", err
		}
		if out == "" || !strings.Contains(out, ns) {
			return "", fmt.Errorf("readlink get invalid: %s", out)
		}

		return out, nil
	}

	for _, ns := range []string{"cgroup", "ipc", "mnt", "net", "pid", "uts"} {
		if _, err := os.Stat("/proc/1/ns/" + ns); err != nil {
			c.Logf("no %s ns, skip", ns)
			continue
		}
		cns, err := getNS(pid, ns)
		c.Assert(err, check.IsNil)
		hns, err := getNS("1", ns)
		c.Assert(err, check.IsNil)

		c.Assert(strings.TrimSpace(cns), check.Not(check.Equals), strings.TrimSpace(hns))
	}
}

// tests runc should set env file in /etc/instanceInfo, we skip test /etc/profile.d/pouchenv.sh
// since it only exist when /etc/profile.d/ exist
func (suite *PouchAliKernelSuite) TestContainerENV(c *check.C) {
	name := "TestContainerENV"
	defer DelContainerForceMultyTime(c, name)

	command.PouchRun("run", "--name", name, "-e", "a=b", "-d", busyboxImage, "top").Assert(c, icmd.Success)

	res := command.PouchRun("exec", name, "env").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "a=b"), check.IsNil)
	res = command.PouchRun("exec", name, "cat", "/etc/instanceInfo").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "env_a = b"), check.IsNil)

	// test env update should take effect in /etc/instanceInfo
	command.PouchRun("update", "-e", "a=newb", name).Assert(c, icmd.Success)
	res = command.PouchRun("exec", name, "env").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "a=newb"), check.IsNil)
	res = command.PouchRun("exec", name, "cat", "/etc/instanceInfo").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Stdout(), "env_a = newb"), check.IsNil)
}

// tests container open file limit should be 655350
// see runc http://gitlab.alibaba-inc.com/pouch/runc/commit/a5e77db5ffedd153cc777491b01b55553a1157bc
func (suite *PouchAliKernelSuite) TestContainerRlimit(c *check.C) {
	name := "TestContainerRlimit"
	defer DelContainerForceMultyTime(c, name)

	command.PouchRun("run", "--name", name, "-d", busyboxImage, "top").Assert(c, icmd.Success)

	pid := strings.TrimSpace(command.PouchRun("inspect", "-f", "{{.State.Pid}}", name).Stdout())
	c.Assert(strings.TrimSpace(pid), check.Not(check.Equals), "0")

	data, err := ioutil.ReadFile("/proc/" + pid + "/limits")
	c.Assert(err, check.IsNil)
	splits := strings.Split(string(data), "\n")

	for _, line := range splits {
		if !strings.Contains(line, "open files") {
			continue
		}

		if !strings.Contains(line, "655350") {
			c.Fatalf("container open files limit should be 655350")
		}
		break
	}
}
