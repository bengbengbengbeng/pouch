package hookplugins

var (
	// AnnotationPrefix represents the prefix of annotation key.
	AnnotationPrefix = "customization."

	// SpecCpusetTrickCpus defines customization.cpuset_trick_cpus
	SpecCpusetTrickCpus = AnnotationPrefix + "cpuset_trick_cpus"
	// SpecCpusetTrickTasks defines customization.cpuset_trick_tasks
	SpecCpusetTrickTasks = AnnotationPrefix + "cpuset_trick_tasks"
	// SpecCpusetTrickExemptTasks defines customization.cpuset_trick_exempt_tasks
	SpecCpusetTrickExemptTasks = AnnotationPrefix + "cpuset_trick_exempt_tasks"
	// SpecCPUBvtWarpNs defines customization.cpu_bvt_warp_ns
	SpecCPUBvtWarpNs = AnnotationPrefix + "cpu_bvt_warp_ns"
	// SpecCpuacctSchedLatSwitch defines customization.cpuacct_sched_lat_switch
	SpecCpuacctSchedLatSwitch = AnnotationPrefix + "cpuacct_sched_lat_switch"

	// SpecMemoryWmarkRatio defines customization.memory_wmark_ratio
	SpecMemoryWmarkRatio = AnnotationPrefix + "memory_wmark_ratio"
	// SpecMemoryExtra defines customization.memory_extra
	SpecMemoryExtra = AnnotationPrefix + "memory_extra"
	// SpecMemoryForceEmptyCtl defines customization.memory_force_empty_ctl
	SpecMemoryForceEmptyCtl = AnnotationPrefix + "memory_force_empty_ctl"
	// SpecMemoryPriority defines customization.memory_priority
	SpecMemoryPriority = AnnotationPrefix + "memory_priority"
	// SpecMemoryUsePriorityOOM defines customization.memory_use_priority_oom
	SpecMemoryUsePriorityOOM = AnnotationPrefix + "memory_use_priority_oom"
	// SpecMemoryOOMKillAll defines customization.memory_oom_kill_all
	SpecMemoryOOMKillAll = AnnotationPrefix + "memory_oom_kill_all"
	// SpecMemoryDroppable defines customization.memory_droppable
	SpecMemoryDroppable = AnnotationPrefix + "memory_droppable"

	// SpecIntelRdtL3Cbm defines customization.intel_rdt_l3_cbm
	SpecIntelRdtL3Cbm = AnnotationPrefix + "intel_rdt_l3_cbm"
	// SpecIntelRdtGroup defines customization.intel_rdt_group
	SpecIntelRdtGroup = AnnotationPrefix + "intel_rdt_group"
	// SpecIntelRdtMba defines customization.intel_rdt_mba
	SpecIntelRdtMba = AnnotationPrefix + "intel_rdt_mba"

	// SpecBlkioFileLevelSwitch defines customization.blkio_file_level_switch
	SpecBlkioFileLevelSwitch = AnnotationPrefix + "blkio_file_level_switch"
	// SpecBlkioBufferWriteBps defines customization.blkio_buffer_write_bps
	SpecBlkioBufferWriteBps = AnnotationPrefix + "blkio_buffer_write_bps"
	// SpecBlkioMetaWriteTps defines customization.blkio_meta_write_tps
	SpecBlkioMetaWriteTps = AnnotationPrefix + "blkio_meta_write_tps"
	// SpecBlkioFileThrottlePath defines customization.blkio_file_throttle_path
	SpecBlkioFileThrottlePath = AnnotationPrefix + "blkio_file_throttle_path"
	// SpecBlkioBufferWriteSwitch defines customization.blkio_buffer_write_switch
	SpecBlkioBufferWriteSwitch = AnnotationPrefix + "blkio_buffer_write_switch"
	// SpecBlkioDeviceBufferWriteBps defines customization.blkio_device_buffer_write_bps
	SpecBlkioDeviceBufferWriteBps = AnnotationPrefix + "blkio_device_buffer_write_bps"
	// SpecBlkioDeviceIdleTime defines customization.blkio_device_idle_time
	SpecBlkioDeviceIdleTime = AnnotationPrefix + "blkio_device_idle_time"
	// SpecBlkioDeviceLatencyTarget defines customization.blkio_device_latency_target
	SpecBlkioDeviceLatencyTarget = AnnotationPrefix + "blkio_device_latency_target"
	// SpecBlkioDeviceReadLowBps defines customization.blkio_device_read_low_bps
	SpecBlkioDeviceReadLowBps = AnnotationPrefix + "blkio_device_read_low_bps"
	// SpecBlkioDeviceReadLowIOps defines customization.blkio_device_read_low_iops
	SpecBlkioDeviceReadLowIOps = AnnotationPrefix + "blkio_device_read_low_iops"
	// SpecBlkioDeviceWriteLowBps defines customization.blkio_device_write_low_bps
	SpecBlkioDeviceWriteLowBps = AnnotationPrefix + "blkio_device_write_low_bps"
	// SpecBlkioDeviceWriteLowIOps defines customization.blkio_device_write_low_iops
	SpecBlkioDeviceWriteLowIOps = AnnotationPrefix + "blkio_device_write_low_iops"

	// SpecNetCgroupRate defines customization.net_cgroup_rate
	SpecNetCgroupRate = AnnotationPrefix + "net_cgroup_rate"
	// SpecNetCgroupCeil defines customization.net_cgroup_ceil
	SpecNetCgroupCeil = AnnotationPrefix + "net_cgroup_ceil"
)

// SupportAnnotation represents the support annotation keys.
var SupportAnnotation = map[string]struct{}{
	SpecCpusetTrickCpus:           {},
	SpecCpusetTrickTasks:          {},
	SpecCpusetTrickExemptTasks:    {},
	SpecCPUBvtWarpNs:              {},
	SpecCpuacctSchedLatSwitch:     {},
	SpecMemoryWmarkRatio:          {},
	SpecMemoryExtra:               {},
	SpecMemoryForceEmptyCtl:       {},
	SpecMemoryPriority:            {},
	SpecMemoryUsePriorityOOM:      {},
	SpecMemoryOOMKillAll:          {},
	SpecMemoryDroppable:           {},
	SpecIntelRdtL3Cbm:             {},
	SpecIntelRdtGroup:             {},
	SpecIntelRdtMba:               {},
	SpecBlkioFileLevelSwitch:      {},
	SpecBlkioBufferWriteBps:       {},
	SpecBlkioMetaWriteTps:         {},
	SpecBlkioFileThrottlePath:     {},
	SpecBlkioBufferWriteSwitch:    {},
	SpecBlkioDeviceBufferWriteBps: {},
	SpecBlkioDeviceIdleTime:       {},
	SpecBlkioDeviceLatencyTarget:  {},
	SpecBlkioDeviceReadLowBps:     {},
	SpecBlkioDeviceReadLowIOps:    {},
	SpecBlkioDeviceWriteLowBps:    {},
	SpecBlkioDeviceWriteLowIOps:   {},
	SpecNetCgroupRate:             {},
	SpecNetCgroupCeil:             {},
}
