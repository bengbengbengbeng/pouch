package apiplugin

import (
	"github.com/alibaba/pouch/apis/types"
)

// InnerContainerJSON is defined for alidocker
type InnerContainerJSON struct {

	// AppArmorProfile are specific for AppArmor to Unix platforms
	AppArmorProfile string `json:"AppArmorProfile,omitempty"`

	// The arguments to the command being run
	Args []string `json:"Args"`

	// config
	Config *types.ContainerConfig `json:"Config,omitempty"`

	// The time the container was created
	Created string `json:"Created,omitempty"`

	// driver
	Driver string `json:"Driver,omitempty"`

	// exec ids of container
	ExecIds []string `json:"ExecIDs"`

	// graph driver
	GraphDriver *types.GraphDriverData `json:"GraphDriver,omitempty"`

	// host config
	HostConfig *InnerHostConfig `json:"HostConfig,omitempty"`

	// The rootfs path of the container on the host.
	HostRootPath string `json:"HostRootPath,omitempty"`

	// the path of container's hostname file on host.
	HostnamePath string `json:"HostnamePath,omitempty"`

	// the path of container's hosts file on host.
	HostsPath string `json:"HostsPath,omitempty"`

	// The ID of the container
	ID string `json:"Id,omitempty"`

	// The container's image
	Image string `json:"Image,omitempty"`

	// the path of container's log file on host.
	LogPath string `json:"LogPath,omitempty"`

	// MountLabel contains the options for the 'mount' command.
	MountLabel string `json:"MountLabel,omitempty"`

	// Set of mount point in a container.
	Mounts []types.MountPoint `json:"Mounts"`

	// name of the created container.
	Name string `json:"Name,omitempty"`

	// NetworkSettings exposes the network settings in the API.
	NetworkSettings *types.NetworkSettings `json:"NetworkSettings,omitempty"`

	// The path to the command being run
	Path string `json:"Path,omitempty"`

	// process label
	ProcessLabel string `json:"ProcessLabel,omitempty"`

	// the path of container's resolvConf file on host.
	ResolvConfPath string `json:"ResolvConfPath,omitempty"`

	// the container's restart time
	RestartCount int64 `json:"RestartCount,omitempty"`

	// The total size of all the files in this container.
	SizeRootFs *int64 `json:"SizeRootFs,omitempty"`

	// The size of files that have been created or changed by this container.
	SizeRw *int64 `json:"SizeRw,omitempty"`

	// snapshotter
	Snapshotter *types.SnapshotterData `json:"Snapshotter,omitempty"`

	// The state of the container.
	State *types.ContainerState `json:"State,omitempty"`
}

// InnerHostConfig is defined for alidocker
type InnerHostConfig struct {

	// Automatically remove the container when the container's process exits. This has no effect if `RestartPolicy` is set.
	AutoRemove bool `json:"AutoRemove,omitempty"`

	// A list of volume bindings for this container. Each volume binding is a string in one of these forms:
	//
	// - `host-src:container-dest` to bind-mount a host path into the container. Both `host-src`, and `container-dest` must be an _absolute_ path.
	// - `host-src:container-dest:ro` to make the bind mount read-only inside the container. Both `host-src`, and `container-dest` must be an _absolute_ path.
	// - `volume-name:container-dest` to bind-mount a volume managed by a volume driver into the container. `container-dest` must be an _absolute_ path.
	// - `volume-name:container-dest:ro` to mount the volume read-only inside the container.  `container-dest` must be an _absolute_ path.
	//
	Binds []string `json:"Binds"`

	// A list of kernel capabilities to add to the container.
	CapAdd []string `json:"CapAdd"`

	// A list of kernel capabilities to drop from the container.
	CapDrop []string `json:"CapDrop"`

	// Cgroup to use for the container.
	Cgroup string `json:"Cgroup,omitempty"`

	// Initial console size, as an `[height, width]` array. (Windows only)
	// Max Items: 2
	// Min Items: 2
	ConsoleSize []*int64 `json:"ConsoleSize"`

	// Path to a file where the container ID is written
	ContainerIDFile string `json:"ContainerIDFile,omitempty"`

	// A list of DNS servers for the container to use.
	DNS []string `json:"Dns"`

	// A list of DNS options.
	DNSOptions []string `json:"DnsOptions"`

	// A list of DNS search domains.
	DNSSearch []string `json:"DnsSearch"`

	// Whether to enable lxcfs.
	EnableLxcfs bool `json:"EnableLxcfs,omitempty"`

	// A list of hostnames/IP mappings to add to the container's `/etc/hosts` file. Specified in the form `["hostname:IP"]`.
	//
	ExtraHosts []string `json:"ExtraHosts"`

	// A list of additional groups that the container process will run as.
	GroupAdd []string `json:"GroupAdd"`

	// Initial script executed in container. The script will be executed before entrypoint or command
	InitScript string `json:"InitScript,omitempty"`

	// IPC sharing mode for the container. Possible values are:
	// - `"none"`: own private IPC namespace, with /dev/shm not mounted
	// - `"private"`: own private IPC namespace
	// - `"shareable"`: own private IPC namespace, with a possibility to share it with other containers
	// - `"container:<name|id>"`: join another (shareable) container's IPC namespace
	// - `"host"`: use the host system's IPC namespace
	// If not specified, daemon default is used, which can either be `"private"`
	// or `"shareable"`, depending on daemon version and configuration.
	//
	IpcMode string `json:"IpcMode,omitempty"`

	// Isolation technology of the container. (Windows only)
	// Enum: [default process hyperv]
	Isolation string `json:"Isolation,omitempty"`

	// A list of links for the container in the form `container_name:alias`.
	Links []string `json:"Links"`

	// The logging configuration for this container
	LogConfig *types.LogConfig `json:"LogConfig,omitempty"`

	// Network mode to use for this container. Supported standard values are: `bridge`, `host`, `none`, and `container:<name|id>`. Any other value is taken as a custom network's name to which this container should connect to.
	NetworkMode string `json:"NetworkMode,omitempty"`

	// An integer value containing the score given to the container in order to tune OOM killer preferences.
	// The range is in [-1000, 1000].
	//
	// Maximum: 1000
	// Minimum: -1000
	OomScoreAdj int64 `json:"OomScoreAdj,omitempty"`

	// Set the PID (Process) Namespace mode for the container. It can be either:
	// - `"container:<name|id>"`: joins another container's PID namespace
	// - `"host"`: use the host's PID namespace inside the container
	//
	PidMode string `json:"PidMode,omitempty"`

	// A map of exposed container ports and the host port they should map to.
	PortBindings types.PortMap `json:"PortBindings,omitempty"`

	// Gives the container full access to the host.
	Privileged bool `json:"Privileged"`

	// Allocates a random host port for all of a container's exposed ports.
	PublishAllPorts bool `json:"PublishAllPorts,omitempty"`

	// Mount the container's root filesystem as read only.
	ReadonlyRootfs bool `json:"ReadonlyRootfs,omitempty"`

	// Restart policy to be used to manage the container
	RestartPolicy *types.RestartPolicy `json:"RestartPolicy,omitempty"`

	// Whether to start container in rich container mode. (default false)
	Rich bool `json:"Rich,omitempty"`

	// Choose one rich container mode.(default dumb-init)
	// Enum: [dumb-init sbin-init systemd]
	RichMode string `json:"RichMode,omitempty"`

	// Runtime to use with this container.
	Runtime string `json:"Runtime,omitempty"`

	// A list of string values to customize labels for MLS systems, such as SELinux.
	SecurityOpt []string `json:"SecurityOpt"`

	// Size of `/dev/shm` in bytes. If omitted, the system uses 64MB.
	// Minimum: 0
	ShmSize *int64 `json:"ShmSize,omitempty"`

	// Storage driver options for this container, in the form `{"size": "120G"}`.
	//
	StorageOpt map[string]string `json:"StorageOpt,omitempty"`

	// A list of kernel parameters (sysctls) to set in the container. For example: `{"net.ipv4.ip_forward": "1"}`
	//
	Sysctls map[string]string `json:"Sysctls,omitempty"`

	// A map of container directories which should be replaced by tmpfs mounts, and their corresponding mount options. For example: `{ "/run": "rw,noexec,nosuid,size=65536k" }`.
	//
	Tmpfs map[string]string `json:"Tmpfs,omitempty"`

	// UTS namespace to use for the container.
	UTSMode string `json:"UTSMode,omitempty"`

	// Sets the usernamespace mode for the container when usernamespace remapping option is enabled.
	UsernsMode string `json:"UsernsMode,omitempty"`

	// Driver that this container uses to mount volumes.
	VolumeDriver string `json:"VolumeDriver,omitempty"`

	// A list of volumes to inherit from another container, specified in the form `<container name>[:<ro|rw>]`.
	VolumesFrom []string `json:"VolumesFrom"`

	InnerResources
}

// InnerResources is defined for alidocker resources
type InnerResources struct {

	// Limit read rate (bytes per second) from a device, in the form `[{"Path": "device_path", "Rate": rate}]`.
	//
	BlkioDeviceReadBps []*types.ThrottleDevice `json:"BlkioDeviceReadBps"`

	// Limit read rate (IO per second) from a device, in the form `[{"Path": "device_path", "Rate": rate}]`.
	//
	BlkioDeviceReadIOps []*types.ThrottleDevice `json:"BlkioDeviceReadIOps"`

	// Limit write rate (bytes per second) to a device, in the form `[{"Path": "device_path", "Rate": rate}]`.
	//
	BlkioDeviceWriteBps []*types.ThrottleDevice `json:"BlkioDeviceWriteBps"`

	// Limit write rate (IO per second) to a device, in the form `[{"Path": "device_path", "Rate": rate}]`.
	//
	BlkioDeviceWriteIOps []*types.ThrottleDevice `json:"BlkioDeviceWriteIOps"`

	// Block IO weight (relative weight), need CFQ IO Scheduler enable.
	// Maximum: 1000
	// Minimum: 0
	BlkioWeight uint16 `json:"BlkioWeight"`

	// Block IO weight (relative device weight) in the form `[{"Path": "device_path", "Weight": weight}]`.
	//
	BlkioWeightDevice []*types.WeightDevice `json:"BlkioWeightDevice"`

	// Path to `cgroups` under which the container's `cgroup` is created. If the path is not absolute, the path is considered to be relative to the `cgroups` path of the init process. Cgroups are created if they do not already exist.
	CgroupParent string `json:"CgroupParent"`

	// The number of usable CPUs (Windows only).
	// On Windows Server containers, the processor resource controls are mutually exclusive. The order of precedence is `CPUCount` first, then `CPUShares`, and `CPUPercent` last.
	//
	CPUCount int64 `json:"CpuCount"`

	// The usable percentage of the available CPUs (Windows only).
	// On Windows Server containers, the processor resource controls are mutually exclusive. The order of precedence is `CPUCount` first, then `CPUShares`, and `CPUPercent` last.
	//
	CPUPercent int64 `json:"CpuPercent"`

	// CPU CFS (Completely Fair Scheduler) period.
	// The length of a CPU period in microseconds.
	//
	// Maximum: 1e+06
	// Minimum: 1000
	CPUPeriod int64 `json:"CpuPeriod"`

	// CPU CFS (Completely Fair Scheduler) quota.
	// Microseconds of CPU time that the container can get in a CPU period."
	//
	CPUQuota int64 `json:"CpuQuota"`

	// The length of a CPU real-time period in microseconds. Set to 0 to allocate no time allocated to real-time tasks.
	CPURealtimePeriod int64 `json:"CpuRealtimePeriod"`

	// The length of a CPU real-time runtime in microseconds. Set to 0 to allocate no time allocated to real-time tasks.
	CPURealtimeRuntime int64 `json:"CpuRealtimeRuntime"`

	// An integer value representing this container's relative CPU weight versus other containers.
	CPUShares int64 `json:"CpuShares"`

	// CPUs in which to allow execution (e.g., `0-3`, `0,1`)
	CpusetCpus string `json:"CpusetCpus"`

	// Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on NUMA systems.
	CpusetMems string `json:"CpusetMems"`

	// a list of cgroup rules to apply to the container
	DeviceCgroupRules []string `json:"DeviceCgroupRules"`

	// A list of devices to add to the container.
	Devices []*types.DeviceMapping `json:"Devices"`

	// Maximum IO in bytes per second for the container system drive (Windows only)
	IOMaximumBandwidth uint64 `json:"IOMaximumBandwidth"`

	// Maximum IOps for the container system drive (Windows only)
	IOMaximumIOps uint64 `json:"IOMaximumIOps"`

	// IntelRdtL3Cbm specifies settings for Intel RDT/CAT group that the container is placed into to limit the resources (e.g., L3 cache) the container has available.
	IntelRdtL3Cbm string `json:"IntelRdtL3Cbm"`

	// Kernel memory limit in bytes.
	KernelMemory int64 `json:"KernelMemory"`

	// Memory limit in bytes.
	Memory int64 `json:"Memory"`

	// MemoryExtra is an integer value representing this container's memory high water mark percentage.
	// The range is in [0, 100].
	//
	// Maximum: 100
	// Minimum: 0
	MemoryExtra *int64 `json:"MemoryExtra"`

	// MemoryForceEmptyCtl represents whether to reclaim the page cache when deleting cgroup.
	// Maximum: 1
	// Minimum: 0
	MemoryForceEmptyCtl int64 `json:"MemoryForceEmptyCtl"`

	// Memory soft limit in bytes.
	MemoryReservation int64 `json:"MemoryReservation"`

	// Total memory limit (memory + swap). Set as `-1` to enable unlimited swap.
	MemorySwap int64 `json:"MemorySwap"`

	// Tune a container's memory swappiness behavior. Accepts an integer between 0 and 100. -1 is also accepted, as a legacy alias of 0.
	// Maximum: 100
	// Minimum: -1
	MemorySwappiness *int64 `json:"MemorySwappiness"`

	// MemoryWmarkRatio is an integer value representing this container's memory low water mark percentage.
	// The value of memory low water mark is memory.limit_in_bytes * MemoryWmarkRatio. The range is in [0, 100].
	//
	// Maximum: 100
	// Minimum: 0
	MemoryWmarkRatio *int64 `json:"MemoryWmarkRatio"`

	// CPU quota in units of 10<sup>-9</sup> CPUs.
	NanoCpus int64 `json:"NanoCpus"`

	// nvidia config
	NvidiaConfig *types.NvidiaConfig `json:"NvidiaConfig,omitempty"`

	// Disable OOM Killer for the container.
	OomKillDisable *bool `json:"OomKillDisable"`

	// Tune a container's pids limit. Set -1 for unlimited. Only on Linux 4.4 does this parameter support.
	//
	PidsLimit int64 `json:"PidsLimit"`

	// ScheLatSwitch enables scheduler latency count in cpuacct
	// Maximum: 1
	// Minimum: 0
	ScheLatSwitch int64 `json:"ScheLatSwitch"`

	// A list of resource limits to set in the container. For example: `{"Name": "nofile", "Soft": 1024, "Hard": 2048}`"
	//
	Ulimits []*types.Ulimit `json:"Ulimits"`

	ResourcesWrapper
}
