package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/test/request"
	"github.com/go-check/check"
)

// TestCreateWithHomeDir is to verify the homedir param of creating container
func (suite *APIContainerCreateSuite) TestCreateWithHomedir(c *check.C) {
	homeDir := "/tmp/homedir/container"
	os.RemoveAll(homeDir)
	defer os.RemoveAll(homeDir)

	cname := "TestCreateWithHomedir"
	q := url.Values{}
	q.Add("name", cname)
	query := request.WithQuery(q)

	obj := map[string]interface{}{
		"Image": busyboxImage,
		"Home":  homeDir,
	}
	body := request.WithJSONBody(obj)

	resp, err := request.Post("/containers/create", query, body)
	defer DelContainerForceMultyTime(c, cname)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)

	if _, err := os.Stat(homeDir); err != nil {
		c.Fatalf("failt to stat home directory(%s) %s", homeDir, err.Error())
	}

	resp, err = request.Get("/containers/" + cname + "/json")
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 200)

	got := types.ContainerJSON{}
	err = request.DecodeBody(&got, resp.Body)
	c.Assert(err, check.IsNil)

	upperDir := got.GraphDriver.Data["UpperDir"]
	if !strings.HasPrefix(upperDir, homeDir) {
		c.Fatalf("expect upperdir(%s) in home dir(%s)", upperDir, homeDir)
	}
}
