package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/hookplugins/apiplugin"
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"

	"github.com/go-check/check"
)

// APIContainerPluginsSuite is the test suite for container API plugin.
type APIContainerPluginsSuite struct{}

func init() {
	check.Suite(&APIContainerPluginsSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APIContainerPluginsSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)

	PullImage(c, busyboxImage)
	PullImage(c, alios7u)
}

// TestCpusetTrick tests creating container cpuset trick.
func (suite *APIContainerPluginsSuite) TestCpusetTrick(c *check.C) {
	q := url.Values{
		"name": []string{"TestCpusetTrick"},
	}
	query := request.WithQuery(q)

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
	resp, err := request.Post("/containers/create", query, body)
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

// DockerContainerCreateConfig is the parameter set to ContainerCreate()
type DockerContainerCreateConfig struct {
	types.ContainerConfig
	HostConfig *apiplugin.InnerHostConfig
}

func createThrottleDevice(rate uint64) []*types.ThrottleDevice {
	tds := []*types.ThrottleDevice{}
	td1 := &types.ThrottleDevice{
		Path: "/tmp",
		Rate: rate,
	}

	td2 := &types.ThrottleDevice{
		Path: "/home",
		Rate: rate,
	}

	tds = append(tds, td1, td2)
	return tds
}

func validateThrottleDevice(src, dst []*types.ThrottleDevice, c *check.C) {
	c.Assert(len(src), check.Equals, len(dst))
	for i := range src {
		c.Assert(src[i].Path, check.Equals, dst[i].Path)
		c.Assert(src[i].Rate, check.Equals, dst[i].Rate)
	}
}

// TestContainerCreateAndInspect tests for alidocker container create and inspect api
func (suite *APIContainerPluginsSuite) TestContainerCreateAndInspect(c *check.C) {
	cname := "api_plugin_container_create_test"

	createConfig := &DockerContainerCreateConfig{
		ContainerConfig: types.ContainerConfig{
			Image: alios7u,
		},
		HostConfig: &apiplugin.InnerHostConfig{
			NetworkMode: "none",
		},
	}

	{
		createConfig.HostConfig.CpusetCpus = "1,2"
		createConfig.HostConfig.CpusetTrickCpus = "1,2"
		createConfig.HostConfig.CpusetTrickTasks = "nginx"
		createConfig.HostConfig.CpusetTrickExemptTasks = "top"
		createConfig.HostConfig.CPUBvtWarpNs = int64(1000)
		createConfig.HostConfig.ScheLatSwitch = int64(1)

		createConfig.HostConfig.Memory = 2 * 1024 * 1024 * 1024

		memoryWmarkRatio := int(50)
		createConfig.HostConfig.MemoryWmarkRatio = memoryWmarkRatio

		memoryExtra := int64(90)
		createConfig.HostConfig.MemoryExtra = memoryExtra
		createConfig.HostConfig.MemoryForceEmptyCtl = int(1)
		memoryPriority := 10
		createConfig.HostConfig.MemoryPriority = &memoryPriority
		memoryUsePriorityOOM := 1
		createConfig.HostConfig.MemoryUsePriorityOOM = &memoryUsePriorityOOM
		createConfig.HostConfig.MemoryDroppable = 1
		memoryKillAll := 1
		createConfig.HostConfig.MemoryKillAll = &memoryKillAll

		createConfig.HostConfig.BlkBufferWriteSwitch = 1
		createConfig.HostConfig.BlkBufferWriteBps = 100000
		createConfig.HostConfig.BlkMetaWriteTps = 1000
		createConfig.HostConfig.BlkFileLevelSwitch = 1
		createConfig.HostConfig.BlkFileThrottlePath = []string{"/tmp"}
		createConfig.HostConfig.BlkioDeviceReadLowIOps = createThrottleDevice(1001)
		createConfig.HostConfig.BlkioDeviceReadLowBps = createThrottleDevice(1002)
		createConfig.HostConfig.BlkioDeviceWriteLowBps = createThrottleDevice(1003)
		createConfig.HostConfig.BlkioDeviceWriteLowIOps = createThrottleDevice(1004)
		createConfig.HostConfig.BlkDeviceBufferWriteBps = createThrottleDevice(1005)
		createConfig.HostConfig.BlkDeviceIdleTime = createThrottleDevice(1006)
		createConfig.HostConfig.BlkDeviceLatencyTarget = createThrottleDevice(1007)

		createConfig.HostConfig.IntelRdtL3Cbm = "L3:0=fffff;1=1ff"
		createConfig.HostConfig.IntelRdtMba = "MB:0=100;1=100"
		createConfig.HostConfig.IntelRdtGroup = "ff11"
	}

	body := request.WithJSONBody(createConfig)
	header := request.WithHeader("User-Agent", "Docker-Client")

	query := make(url.Values)
	query.Add("name", cname)
	params := request.WithQuery(query)

	resp, err := request.Post("/containers/create", body, header, params)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)

	defer DelContainerForceMultyTime(c, cname)

	resp, err = request.Get(fmt.Sprintf("/containers/%s/json", cname), header)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 200)

	inspectJSON := &apiplugin.InnerContainerJSON{}
	err = json.NewDecoder(resp.Body).Decode(inspectJSON)
	c.Assert(err, check.IsNil)

	c.Assert(inspectJSON.HostConfig.NetworkMode, check.Equals, "none")
	c.Assert(inspectJSON.HostConfig.CpusetCpus, check.Equals, "1,2")
	c.Assert(inspectJSON.HostConfig.CpusetTrickCpus, check.Equals, "1,2")
	c.Assert(inspectJSON.HostConfig.CpusetTrickTasks, check.Equals, "nginx")
	c.Assert(inspectJSON.HostConfig.CpusetTrickExemptTasks, check.Equals, "top")
	c.Assert(inspectJSON.HostConfig.CPUBvtWarpNs, check.Equals, int64(1000))
	c.Assert(inspectJSON.HostConfig.ScheLatSwitch, check.Equals, int64(1))

	c.Assert(inspectJSON.HostConfig.Memory, check.Equals, int64(2*1024*1024*1024))
	c.Assert(inspectJSON.HostConfig.MemoryWmarkRatio, check.Equals, int(50))
	c.Assert(inspectJSON.HostConfig.MemoryExtra, check.Equals, int64(90))
	c.Assert(inspectJSON.HostConfig.MemoryForceEmptyCtl, check.Equals, int(1))
	c.Assert(*inspectJSON.HostConfig.MemoryPriority, check.Equals, 10)
	c.Assert(*inspectJSON.HostConfig.MemoryUsePriorityOOM, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.MemoryDroppable, check.Equals, 1)
	c.Assert(*inspectJSON.HostConfig.MemoryKillAll, check.Equals, 1)

	c.Assert(inspectJSON.HostConfig.BlkBufferWriteSwitch, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.BlkBufferWriteBps, check.Equals, 100000)
	c.Assert(inspectJSON.HostConfig.BlkMetaWriteTps, check.Equals, 1000)
	c.Assert(inspectJSON.HostConfig.BlkFileLevelSwitch, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.BlkFileThrottlePath, check.DeepEquals, []string{"/tmp"})
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceReadLowIOps, createThrottleDevice(1001), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceReadLowBps, createThrottleDevice(1002), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceWriteLowBps, createThrottleDevice(1003), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceWriteLowIOps, createThrottleDevice(1004), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceBufferWriteBps, createThrottleDevice(1005), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceIdleTime, createThrottleDevice(1006), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceLatencyTarget, createThrottleDevice(1007), c)

	c.Assert(inspectJSON.HostConfig.IntelRdtL3Cbm, check.Equals, "L3:0=fffff;1=1ff")
	c.Assert(inspectJSON.HostConfig.IntelRdtMba, check.Equals, "MB:0=100;1=100")
	c.Assert(inspectJSON.HostConfig.IntelRdtGroup, check.Equals, "ff11")
}

// DockerContainerUpdateConfig is the parameter set to ContainerUpdate api
type DockerContainerUpdateConfig struct {
	apiplugin.ResourcesWrapper
}

// TestContainerUpdateForApiPlugin tests for alidocker container update api
func (suite *APIContainerPluginsSuite) TestContainerUpdateForApiPlugin(c *check.C) {
	cname := "api_plugin_container_update_test"

	createConfig := &DockerContainerCreateConfig{
		ContainerConfig: types.ContainerConfig{
			Image: alios7u,
		},
		HostConfig: &apiplugin.InnerHostConfig{
			NetworkMode: "none",
		},
	}

	{
		createConfig.HostConfig.CpusetCpus = "1,2"
		createConfig.HostConfig.CpusetTrickCpus = "1,2"
		createConfig.HostConfig.CpusetTrickTasks = "nginx"
		createConfig.HostConfig.CpusetTrickExemptTasks = "top"
		createConfig.HostConfig.CPUBvtWarpNs = int64(1000)
		createConfig.HostConfig.ScheLatSwitch = int64(1)

		createConfig.HostConfig.Memory = 2 * 1024 * 1024 * 1024

		memoryWmarkRatio := int(50)
		createConfig.HostConfig.MemoryWmarkRatio = memoryWmarkRatio

		memoryExtra := int64(90)
		createConfig.HostConfig.MemoryExtra = memoryExtra
		createConfig.HostConfig.MemoryForceEmptyCtl = int(1)
		memoryPriority := 10
		createConfig.HostConfig.MemoryPriority = &memoryPriority
		memoryUsePriorityOOM := 1
		createConfig.HostConfig.MemoryUsePriorityOOM = &memoryUsePriorityOOM
		createConfig.HostConfig.MemoryDroppable = 1
		memoryKillAll := 1
		createConfig.HostConfig.MemoryKillAll = &memoryKillAll

		createConfig.HostConfig.BlkBufferWriteSwitch = 1
		createConfig.HostConfig.BlkBufferWriteBps = 100000
		createConfig.HostConfig.BlkMetaWriteTps = 1000
		createConfig.HostConfig.BlkFileLevelSwitch = 1
		createConfig.HostConfig.BlkFileThrottlePath = []string{"/tmp"}
		createConfig.HostConfig.BlkioDeviceReadLowIOps = createThrottleDevice(1001)
		createConfig.HostConfig.BlkioDeviceReadLowBps = createThrottleDevice(1002)
		createConfig.HostConfig.BlkioDeviceWriteLowBps = createThrottleDevice(1003)
		createConfig.HostConfig.BlkioDeviceWriteLowIOps = createThrottleDevice(1004)
		createConfig.HostConfig.BlkDeviceBufferWriteBps = createThrottleDevice(1005)
		createConfig.HostConfig.BlkDeviceIdleTime = createThrottleDevice(1006)
		createConfig.HostConfig.BlkDeviceLatencyTarget = createThrottleDevice(1007)

		createConfig.HostConfig.IntelRdtL3Cbm = "L3:0=fffff;1=1ff"
		createConfig.HostConfig.IntelRdtMba = "MB:0=100;1=100"
		createConfig.HostConfig.IntelRdtGroup = "ff11"
	}

	body := request.WithJSONBody(createConfig)
	header := request.WithHeader("User-Agent", "Docker-Client")

	query := make(url.Values)
	query.Add("name", cname)
	params := request.WithQuery(query)

	resp, err := request.Post("/containers/create", body, header, params)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 201)

	defer DelContainerForceMultyTime(c, cname)

	updateConfig := &DockerContainerUpdateConfig{}

	{
		updateConfig.ResourcesWrapper.CpusetTrickCpus = "1"
		updateConfig.ResourcesWrapper.CpusetTrickTasks = "top"
		updateConfig.ResourcesWrapper.CpusetTrickExemptTasks = "nginx"
		updateConfig.ResourcesWrapper.CPUBvtWarpNs = int64(2)
		updateConfig.ResourcesWrapper.ScheLatSwitch = int64(1)

		memoryWmarkRatio := int(80)
		updateConfig.ResourcesWrapper.MemoryWmarkRatio = memoryWmarkRatio

		updateConfig.ResourcesWrapper.MemoryExtra = int64(80)
		memoryPriority := 9
		updateConfig.ResourcesWrapper.MemoryPriority = &memoryPriority
		// todo: memoryUsePriorityOOM could be set to 0 and MemoryKillAll could be set to 0?
		// memoryUsePriorityOOM := 0
		// updateConfig.ResourcesWrapper.MemoryUsePriorityOOM = &memoryUsePriorityOOM
		// memoryKillAll := 0
		// updateConfig.ResourcesWrapper.MemoryKillAll = &memoryKillAll

		updateConfig.ResourcesWrapper.BlkBufferWriteBps = 200000
		updateConfig.ResourcesWrapper.BlkMetaWriteTps = 2000
		updateConfig.ResourcesWrapper.BlkFileThrottlePath = []string{"/home"}
		updateConfig.ResourcesWrapper.BlkioDeviceReadLowIOps = createThrottleDevice(101)
		updateConfig.ResourcesWrapper.BlkioDeviceReadLowBps = createThrottleDevice(102)
		updateConfig.ResourcesWrapper.BlkioDeviceWriteLowBps = createThrottleDevice(103)
		updateConfig.ResourcesWrapper.BlkioDeviceWriteLowIOps = createThrottleDevice(104)
		updateConfig.ResourcesWrapper.BlkDeviceBufferWriteBps = createThrottleDevice(105)
		updateConfig.ResourcesWrapper.BlkDeviceIdleTime = createThrottleDevice(106)
		updateConfig.ResourcesWrapper.BlkDeviceLatencyTarget = createThrottleDevice(107)

		updateConfig.ResourcesWrapper.IntelRdtL3Cbm = "L3:0=fffff;1=1fe"
		updateConfig.ResourcesWrapper.IntelRdtMba = "MB:0=99;1=100"
		updateConfig.ResourcesWrapper.IntelRdtGroup = "ff22"
	}

	updateBody := request.WithJSONBody(updateConfig)
	resp, err = request.Post(fmt.Sprintf("/containers/%s/update", cname), updateBody, header)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 200)

	resp, err = request.Get(fmt.Sprintf("/containers/%s/json", cname), header)
	c.Assert(err, check.IsNil)
	CheckRespStatus(c, resp, 200)

	inspectJSON := &apiplugin.InnerContainerJSON{}
	err = json.NewDecoder(resp.Body).Decode(inspectJSON)
	c.Assert(err, check.IsNil)

	c.Assert(inspectJSON.HostConfig.NetworkMode, check.Equals, "none")
	c.Assert(inspectJSON.HostConfig.CpusetCpus, check.Equals, "1,2")
	c.Assert(inspectJSON.HostConfig.CpusetTrickCpus, check.Equals, "1")
	c.Assert(inspectJSON.HostConfig.CpusetTrickTasks, check.Equals, "top")
	c.Assert(inspectJSON.HostConfig.CpusetTrickExemptTasks, check.Equals, "nginx")
	c.Assert(inspectJSON.HostConfig.CPUBvtWarpNs, check.Equals, int64(2))
	c.Assert(inspectJSON.HostConfig.ScheLatSwitch, check.Equals, int64(1))

	c.Assert(inspectJSON.HostConfig.Memory, check.Equals, int64(2*1024*1024*1024))
	c.Assert(inspectJSON.HostConfig.MemoryWmarkRatio, check.Equals, int(80))
	c.Assert(inspectJSON.HostConfig.MemoryExtra, check.Equals, int64(80))
	c.Assert(inspectJSON.HostConfig.MemoryForceEmptyCtl, check.Equals, int(1))
	c.Assert(*inspectJSON.HostConfig.MemoryPriority, check.Equals, 9)
	c.Assert(*inspectJSON.HostConfig.MemoryUsePriorityOOM, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.MemoryDroppable, check.Equals, 1)
	c.Assert(*inspectJSON.HostConfig.MemoryKillAll, check.Equals, 1)

	c.Assert(inspectJSON.HostConfig.BlkBufferWriteSwitch, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.BlkBufferWriteBps, check.Equals, 200000)
	c.Assert(inspectJSON.HostConfig.BlkMetaWriteTps, check.Equals, 2000)
	c.Assert(inspectJSON.HostConfig.BlkFileLevelSwitch, check.Equals, 1)
	c.Assert(inspectJSON.HostConfig.BlkFileThrottlePath, check.DeepEquals, []string{"/home"})
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceReadLowIOps, createThrottleDevice(101), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceReadLowBps, createThrottleDevice(102), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceWriteLowBps, createThrottleDevice(103), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkioDeviceWriteLowIOps, createThrottleDevice(104), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceBufferWriteBps, createThrottleDevice(105), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceIdleTime, createThrottleDevice(106), c)
	validateThrottleDevice(inspectJSON.HostConfig.BlkDeviceLatencyTarget, createThrottleDevice(107), c)

	c.Assert(inspectJSON.HostConfig.IntelRdtL3Cbm, check.Equals, "L3:0=fffff;1=1fe")
	c.Assert(inspectJSON.HostConfig.IntelRdtMba, check.Equals, "MB:0=99;1=100")
	c.Assert(inspectJSON.HostConfig.IntelRdtGroup, check.Equals, "ff22")
}
