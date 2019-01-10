package main

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
)

// APIContainerCreateSuite is the test suite for container create API.
type APIContainerPluginsSuite struct{}

func init() {
	check.Suite(&APIContainerPluginsSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APIContainerPluginsSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)

	PullImage(c, busyboxImage)
}

// TestCpusetTrick tests creating container cpuset trick.
func (suite *APIContainerPluginsSuite) TestCpusetTrick(c *check.C) {
	obj := map[string]interface{}{
		"Cmd":   []string{"top"},
		"Image": busyboxImage,
		"HostConfig": map[string]interface{}{
			"CpusetCpus":             "0-3",
			"CpusetTrickCpus":        "0,2",
			"CpusetTrickTasks":       "java,nginx",
			"CpusetTrickExemptTasks": "top",
		},
	}

	body := request.WithJSONBody(obj)
	resp, err := request.Post("/containers/create", body)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)

	// Decode response
	got := types.ContainerCreateResp{}
	c.Assert(request.DecodeBody(&got, resp.Body), check.IsNil)
	c.Assert(got.ID, check.NotNil)
	c.Assert(got.Name, check.NotNil)

	defer DelContainerForceMultyTime(c, got.Name)

	resp, err = request.Post("/containers/" + got.Name + "/start")
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 204)

	cgroupCpusetCpus := "/sys/fs/cgroup/cpu/default/" + got.ID + "/cpuset.cpus"
	if _, err := os.Stat(cgroupCpusetCpus); err != nil && !os.IsExist(err) {
		c.Fatalf("cpuset.cpus cgroup is not exist.")
	}
	cgroupCpusetTrickCpus := "/sys/fs/cgroup/cpu/default/" + got.ID + "/cpuset.trick_cpus"
	if _, err := os.Stat(cgroupCpusetTrickCpus); err != nil && !os.IsExist(err) {
		c.Skip("skip test, kernel is not support cpuset.trick_cpus")
	}
	cgroupCpusetTrickTasks := "/sys/fs/cgroup/cpu/default/" + got.ID + "/cpuset.trick_tasks"
	if _, err := os.Stat(cgroupCpusetTrickTasks); err != nil && !os.IsExist(err) {
		c.Skip("skip test, kernel is not support cpuset.trick_tasks")
	}
	cgroupCpusetTrickExemptTasks := "/sys/fs/cgroup/cpu/default/" + got.ID + "/cpuset.trick_exempt_tasks"
	if _, err := os.Stat(cgroupCpusetTrickExemptTasks); err != nil && !os.IsExist(err) {
		c.Skip("skip test, kernel is not support cpuset.trick_exempt_tasks")
	}

	content, err := ioutil.ReadFile(cgroupCpusetCpus)
	if err != nil {
		c.Fatalf("failed to read file: %s", cgroupCpusetCpus)
	}
	if !bytes.Contains(content, []byte("0-3")) {
		c.Fatalf("invalid content: %s, expect [0-3]", content)
	}

	content, err = ioutil.ReadFile(cgroupCpusetTrickCpus)
	if err != nil {
		c.Fatalf("failed to read file: %s", cgroupCpusetTrickCpus)
	}
	if !bytes.Contains(content, []byte("0,2")) {
		c.Fatalf("invalid content: %s, expect [0,2]", content)
	}

	content, err = ioutil.ReadFile(cgroupCpusetTrickTasks)
	if err != nil {
		c.Fatalf("failed to read file: %s", cgroupCpusetTrickTasks)
	}
	if !bytes.Contains(content, []byte("java,nginx")) {
		c.Fatalf("invalid content: %s, expect [java,nginx]", content)
	}

	content, err = ioutil.ReadFile(cgroupCpusetTrickExemptTasks)
	if err != nil {
		c.Fatalf("failed to read file: %s", cgroupCpusetTrickExemptTasks)
	}
	if !bytes.Contains(content, []byte("top")) {
		c.Fatalf("invalid content: %s, expect [top]", content)
	}
}
