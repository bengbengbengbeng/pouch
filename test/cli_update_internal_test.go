package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	hp "github.com/alibaba/pouch/hookplugins"
	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// PouchUpdateInternalSuite is the test suite for internal update CLI.
type PouchUpdateInternalSuite struct{}

func init() {
	check.Suite(&PouchUpdateInternalSuite{})
}

const (
	cgroupBasePath = "/sys/fs/cgroup"
)

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchUpdateInternalSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)

	environment.PruneAllContainers(apiClient)

	PullImage(c, busyboxImage)
}

// TearDownTest does cleanup work in the end of each test.
func (suite *PouchUpdateInternalSuite) TearDownTest(c *check.C) {
}

func checkCgroupExist(cgroupSys string, cgroupFile string) (bool, error) {
	path := filepath.Join(cgroupBasePath, cgroupSys, cgroupFile)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func readCgroupValue(cid string, cgroupSys string, cgroupFile string) (string, error) {
	path := filepath.Join(cgroupBasePath, cgroupSys, "default", cid, cgroupFile)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// TestUpdateCpuBvtWarpNs is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateCpuBvtWarpNs(c *check.C) {
	cgroupExist, err := checkCgroupExist("cpu", "cpu.bvt_warp_ns")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "cpu.bvt_warp_ns"))
	}

	cname := "TestUpdateCpuBvtWarpNs"
	originAnnotation := fmt.Sprintf("%s=1", hp.SpecCPUBvtWarpNs)
	updateAnnotation := fmt.Sprintf("%s=2", hp.SpecCPUBvtWarpNs)

	command.PouchRun("run", "-d", "--name", cname, "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)
	checkContainerAnnotation(c, cname, hp.SpecCPUBvtWarpNs, "2")
	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	value, err := readCgroupValue(cid, "cpu", "cpu.bvt_warp_ns")
	c.Assert(err, check.IsNil)

	c.Check(value, check.Equals, "2")
}

// TestUpdateMemoryForceEmptyCtl is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryForceEmptyCtl(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.force_empty_ctl")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.force_empty_ctl"))
	}

	cname := "TestUpdateMemoryForceEmptyCtl"
	originAnnotation := fmt.Sprintf("%s=1", hp.SpecMemoryForceEmptyCtl)
	updateAnnotation := fmt.Sprintf("%s=0", hp.SpecMemoryForceEmptyCtl)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryForceEmptyCtl, "1")
	value, err := readCgroupValue(cid, "memory", "memory.force_empty_ctl")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "1")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryForceEmptyCtl, "0")
	value, err = readCgroupValue(cid, "memory", "memory.force_empty_ctl")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "0")
}

// TestUpdateMemoryUsePriorityOOM is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryUsePriorityOOM(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.use_priority_oom")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.use_priority_oom"))
	}

	cname := "TestUpdateMemoryUsePriorityOOM"
	originAnnotation := fmt.Sprintf("%s=1", hp.SpecMemoryUsePriorityOOM)
	updateAnnotation := fmt.Sprintf("%s=0", hp.SpecMemoryUsePriorityOOM)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryUsePriorityOOM, "1")
	value, err := readCgroupValue(cid, "memory", "memory.use_priority_oom")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "1")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryUsePriorityOOM, "0")
	value, err = readCgroupValue(cid, "memory", "memory.use_priority_oom")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "0")
}

// TestUpdateMemoryPriority is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryPriority(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.priority")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.priority"))
	}

	cname := "TestUpdateMemoryPriority"
	originAnnotation := fmt.Sprintf("%s=11", hp.SpecMemoryPriority)
	updateAnnotation := fmt.Sprintf("%s=2", hp.SpecMemoryPriority)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryPriority, "11")
	value, err := readCgroupValue(cid, "memory", "memory.priority")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "11")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryPriority, "2")
	value, err = readCgroupValue(cid, "memory", "memory.priority")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "2")
}

// TestUpdateMemoryOOMKillAll is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryOOMKillAll(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.oom_kill_all")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.oom_kill_all"))
	}

	cname := "TestUpdateMemoryOOMKillAll"
	originAnnotation := fmt.Sprintf("%s=0", hp.SpecMemoryOOMKillAll)
	updateAnnotation := fmt.Sprintf("%s=1", hp.SpecMemoryOOMKillAll)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryOOMKillAll, "0")
	value, err := readCgroupValue(cid, "memory", "memory.oom_kill_all")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "0")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryOOMKillAll, "1")
	value, err = readCgroupValue(cid, "memory", "memory.oom_kill_all")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "1")
}

// TestUpdateMemoryWmarkRatio is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryWmarkRatio(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.wmark_ratio")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.wmark_ratio"))
	}

	cname := "TestUpdateMemoryOOMKillAll"
	originAnnotation := fmt.Sprintf("%s=99", hp.SpecMemoryWmarkRatio)
	updateAnnotation := fmt.Sprintf("%s=1", hp.SpecMemoryWmarkRatio)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryWmarkRatio, "99")
	value, err := readCgroupValue(cid, "memory", "memory.wmark_ratio")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "99")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryWmarkRatio, "1")
	value, err = readCgroupValue(cid, "memory", "memory.wmark_ratio")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "1")
}

// TestUpdateMemoryWmarkRatio is to verity the correctness of update the annotation
func (suite *PouchUpdateInternalSuite) TestUpdateMemoryDroppable(c *check.C) {
	cgroupExist, err := checkCgroupExist("memory", "memory.droppable")
	c.Assert(err, check.IsNil)

	if !cgroupExist {
		c.Skip(fmt.Sprintf("cgroup %s is not exist", "memory.droppable"))
	}

	cname := "TestUpdateMemoryOOMKillAll"
	originAnnotation := fmt.Sprintf("%s=1", hp.SpecMemoryDroppable)
	updateAnnotation := fmt.Sprintf("%s=0", hp.SpecMemoryDroppable)

	command.PouchRun("run", "-d", "--name", cname, "--memory", "1g", "--annotation", originAnnotation, busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, cname)

	cid, err := inspectFilter(cname, ".ID")
	c.Assert(err, check.IsNil)

	checkContainerAnnotation(c, cname, hp.SpecMemoryDroppable, "1")
	value, err := readCgroupValue(cid, "memory", "memory.droppable")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "1")

	command.PouchRun("update", "--annotation", updateAnnotation, cname).Assert(c, icmd.Success)

	checkContainerAnnotation(c, cname, hp.SpecMemoryDroppable, "0")
	value, err = readCgroupValue(cid, "memory", "memory.droppable")
	c.Assert(err, check.IsNil)
	c.Check(value, check.Equals, "0")
}
