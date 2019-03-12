package main

import (
	"github.com/alibaba/pouch/test/command"
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/util"

	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// PouchKataSuite is the test suite for run CLI.
type PouchKataSuite struct{}

func init() {
	check.Suite(&PouchKataSuite{})
}

// SetUpSuite does common setup in the beginning of each test suite.
func (suite *PouchKataSuite) SetUpSuite(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
	SkipIfFalse(c, environment.IsPhysicalHost)

	environment.PruneAllContainers(apiClient)

	PullImage(c, busyboxImage)
}

// TestRunKata is to verify the correctness of run container with specified name.
func (suite *PouchKataSuite) TestRunKata(c *check.C) {
	name := "TestRunKata"

	command.PouchRun("run", "-d", "--name", name, "--net", "none",
		"--runtime", "kata-runtime",
		busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, name)

	res := command.PouchRun("ps").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Combined(), name), check.IsNil)
}

// TestRunKataPrintHi is to verify run container with executing a command.
func (suite *PouchKataSuite) TestRunKataPrintHi(c *check.C) {
	name := "TestRunKataPrintHi"

	res := command.PouchRun("run", "--name", name, "--net", "none",
		"--runtime", "kata-runtime",
		busyboxImage, "echo", "hi").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, name)

	c.Assert(util.PartialEqual(res.Combined(), "hi"), check.IsNil)
}

// TestKataExec tests exec kata container successful
func (suite *PouchKataSuite) TestKataExec(c *check.C) {
	name := "TestKataExec"
	command.PouchRun("run", "-d", "--name", name, "--net", "none",
		"--runtime", "kata-runtime",
		busyboxImage, "top").Assert(c, icmd.Success)
	defer DelContainerForceMultyTime(c, name)

	res := command.PouchRun("exec", name, "ls").Assert(c, icmd.Success)
	c.Assert(util.PartialEqual(res.Combined(), "etc"), check.IsNil)
}
