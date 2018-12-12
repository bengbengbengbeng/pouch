package main

import (
	"strconv"
	"strings"
	"time"

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

// TestDefaultNetworkMode: verify that default network mode is bridge
func (suite *PouchRunSuite) TestDefaultNetworkMode(c *check.C) {
	name := "TestDefaultNetworkMode"
	res := command.PouchRun("run", "-d", "--name", name, busyboxImage, "top")
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)
	networkmode, err := inspectFilter(name, ".HostConfig.NetworkMode")
	c.Assert(err, check.IsNil)
	c.Assert(networkmode, check.Equals, "bridge")
}

// TestUidFromIp: generate admin uid if env ali_admin_uid=0 exist
// uid=500 + ip with ending number, for example, 10.44.55.67 => uid 567
func (suite *PouchRunSuite) TestUidFromIp(c *check.C) {
	name := "TestUidFromIp"
	endingnumber := 13 // uid=500+13=513
	FakeIP := "192.168.5." + strconv.Itoa(endingnumber)
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=vm", "-e", "RequestedIP="+FakeIP, "--env", "ali_admin_uid=0", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	//time.Sleep(30 * time.Second)

	expectedstring := strconv.Itoa(500 + endingnumber)
	cmd := "id -u admin"
	out := command.PouchRun("exec", "-i", name, "bash", "-c", cmd).Stdout()

	if strings.Compare(strings.Replace(out, "\n", "", -1), expectedstring) != 0 {
		c.Errorf("%s should be equal to %s", out, expectedstring)
	}
}

//TestSetUserToRootInRichContainer: set user to root if running in rich container mode
func (suite *PouchRunSuite) TestSetUserToRootInRichContainer(c *check.C) {
	name := "TestSetUserToRootInRichContainer"
	res := command.PouchRun("run", "-d", "-u", "admin", "--env", "ali_run_mode=vm", "--env", "ali_admin_uid=0", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	//time.Sleep(60 * time.Second)
	expectedstring := "root"
	cmd := "ls -l /tmp/entry.log"
	out := ""

	i := 0
	for i < 20 {
		out = command.PouchRun("exec", "-i", name, "bash", "-c", cmd).Stdout()
		if !strings.Contains(out, expectedstring) {
			i++
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if i >= 20 {
		c.Errorf("%s should contains %s", out, expectedstring)
	}
}

// TestConvertDiskQuotaLabel: convert label DiskQuota to DiskQuota in ContainerConfig parameter
func (suite *PouchRunSuite) TestConvertDiskQuotaLabel(c *check.C) {
	name := "TestConvertDiskQuotaLabel"
	res := command.PouchRun("run", "-d", "-l", "DiskQuota=\"/=1G\"", "--env", "ali_run_mode=vm", "--env", "ali_admin_uid=0", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "1G"
	output := command.PouchRun("inspect", "-f", "{{.Config.DiskQuota}}", name).Stdout()
	if !strings.Contains(output, expectedstring) {
		c.Errorf("%s should contains %s", output, expectedstring)
	}
}

//TestAliRunModeWithCommonVm: in rich container mode, change ali_run_mode=common_vm to ali_run_mode=vm
func (suite *PouchRunSuite) TestAliRunModeWithCommonVm(c *check.C) {
	name := "TestAliRunModeWithCommonVm"
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--env", "ali_admin_uid=0", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "ali_run_mode=vm"
	output := command.PouchRun("inspect", "-f", "{{.Config.Env}}", name).Stdout()
	if !strings.Contains(output, expectedstring) {
		c.Errorf("%s should contains %s", output, expectedstring)
	}
}

//TestLabelsToEnv: in rich container mode, change ali_run_mode=common_vm to ali_run_mode=vm
func (suite *PouchRunSuite) TestLabelsToEnv(c *check.C) {
	name := "TestLabelsToEnv"
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--env", "ali_admin_uid=0", "--env", "RequestedIP=192.168.5.11", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "ali_run_mode=\"vm\""
	cmd := "cat /etc/profile.d/pouchenv.sh"
	out := command.PouchRun("exec", name, "bash", "-c", cmd).Stdout()
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}

	expectedstring = "ali_admin_uid=\"511\""
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}

	expectedstring = "RequestedIP=\"192.168.5.11\""
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}
}

// TestCapabilitiesInRichContainer: in rich container mode, add some capabilities by default
// SYS_RESOURCE SYS_MODULE SYS_PTRACE SYS_PACCT NET_ADMIN SYS_ADMIN are added.
func (suite *PouchRunSuite) TestCapabilitiesInRichContainer(c *check.C) {
	name := "TestCapabilitiesInRichContainer"
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--env", "ali_admin_uid=0", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "SYS_RESOURCE SYS_MODULE SYS_PTRACE SYS_PACCT NET_ADMIN SYS_ADMIN"
	output := command.PouchRun("inspect", "-f", "{{.HostConfig.CapAdd}}", name).Stdout()

	if !strings.Contains(output, expectedstring) {
		c.Errorf("%s should contains %s", output, expectedstring)
	}
}

// TestBindHostsHostnameResolvInRichContainer: in rich container mode, bind /etc/hosts /etc/hostname /etc/resolv.conf files into container
func (suite *PouchRunSuite) TestBindHostsHostnameResolvInRichContainer(c *check.C) {
	name := "TestBindHostsHostnameResolvInRichContainer"
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--env", "ali_admin_uid=0", "-v", "/etc/:/tmp/etc/", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	cmd := "diff /etc/resolv.conf /tmp/etc/resolv.conf"
	command.PouchRun("exec", name, "bash", "-c", cmd).Assert(c, icmd.Success)

	//TODO: /etc/hostname /etc/hosts
}

// TestShmSizeIsHalfOfMemory: in rich container mode, set ShmSize to half of the limit of memory
func (suite *PouchRunSuite) TestShmSizeIsHalfOfMemory(c *check.C) {
	name := "TestShmSizeIsHalfOfMemory"
	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--memory=8G", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)
	output := command.PouchRun("inspect", "-f", "{{.HostConfig.ShmSize}}", name).Stdout()
	c.Assert(strings.TrimSpace(output), check.Equals, "4294967296") // 4294967296=8x1024x1024x1024/2
}

// TestSetHostnameEnv: set HOSTNAME env if HostName specified
func (suite *PouchRunSuite) TestSetHostnameEnv(c *check.C) {
	name := "TestSetHostnameEnv"

	res := command.PouchRun("run", "-d", "--env", "ali_run_mode=common_vm", "--env", "HOSTNAME=myhello", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "HOSTNAME=\"myhello\""
	cmd := "cat /etc/profile.d/pouchenv.sh"
	out := command.PouchRun("exec", name, "bash", "-c", cmd).Stdout()
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}
}

// TestTrimPrefixContainerSlash: if VolumesFrom specified and the container name has a prefix of slash, trim it
func (suite *PouchRunSuite) TestTrimPrefixContainerSlash(c *check.C) {
	con1name := "TestTrimPrefixContainerSlashcon1"
	con2name := "TestTrimPrefixContainerSlashcon2"

	vol1 := "myvol1"
	vol2 := "myvol2"
	command.PouchRun("volume", "create", "--name", vol1).Assert(c, icmd.Success)
	command.PouchRun("volume", "create", "--name", vol2).Assert(c, icmd.Success)
	defer command.PouchRun("volume", "remove", vol1)
	defer command.PouchRun("volume", "remove", vol2)

	con1id := command.PouchRun("run", "-d", "-v", vol1+":/v1", "-v", vol2+":/v2", "--name", con1name, Image7u).Stdout()
	defer DelContainerForceMultyTime(c, con1name)

	cmd := "echo hellopouch > /v1/tmpfile"
	command.PouchRun("exec", con1name, "bash", "-c", cmd).Assert(c, icmd.Success)
	command.PouchRun("stop", con1name).Assert(c, icmd.Success)

	command.PouchRun("run", "-d", "--volumes-from", "/"+con1id, "--name", con2name, Image7u).Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, con2name)
	cmd = "cat /v1/tmpfile"
	//time.Sleep(20 * time.Second)
	out := command.PouchRun("exec", con2name, "bash", "-c", cmd).Stdout()

	expectedstring := "hellopouch"
	if !strings.Contains(out, expectedstring) {
		c.Errorf("%s should contains %s", out, expectedstring)
	}
}

// TestNetPriority: add net-priority into spec-annotations
func (suite *PouchRunSuite) TestNetPriority(c *check.C) {
	name := "TestNetPriority"

	res := command.PouchRun("run", "-d", "--net-priority=7", "--name", name, Image7u)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	expectedstring := "net-priority:7"
	output := command.PouchRun("inspect", "-f", "{{.Config.SpecAnnotation}}", name).Stdout()
	if !strings.Contains(output, expectedstring) {
		c.Errorf("%s should contains %s", output, expectedstring)
	}

}
