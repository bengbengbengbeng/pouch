package main

import (
	"os"
	"strings"

	"github.com/alibaba/pouch/test/command"
	"github.com/go-check/check"
	"github.com/gotestyourself/gotestyourself/icmd"
)

// TestCreateWithHomeDir is to verify the homedir param of creating container
func (suite *PouchCreateSuite) TestCreateWithHomeDir(c *check.C) {
	homeDir := "/tmp/homedir/container"
	os.RemoveAll(homeDir)
	defer os.RemoveAll(homeDir)
	name := "create-with-homedir"
	res := command.PouchRun("create", "--name", name, "--home", homeDir, busyboxImage)
	defer DelContainerForceMultyTime(c, name)
	res.Assert(c, icmd.Success)

	if _, err := os.Stat(homeDir); err != nil {
		c.Fatalf("failt to stat home directory(%s) %s", homeDir, err.Error())
	}

	upperDir, err := inspectFilter(name, ".GraphDriver.Data.UpperDir")
	c.Assert(err, check.IsNil)
	if !strings.HasPrefix(upperDir, homeDir) {
		c.Fatalf("expect upperdir(%s) in home dir(%s)", upperDir, homeDir)
	}
}
