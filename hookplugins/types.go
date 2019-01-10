package hookplugins

var (
	// AnnotationPrefix represents the prefix of annotation key.
	AnnotationPrefix              = "customization."
	SpecCpusetTrickCpus           = AnnotationPrefix + "cpuset_trick_cpus"
	SpecCpusetTrickTasks          = AnnotationPrefix + "cpuset_trick_tasks"
	SpecCpusetTrickExemptTasks    = AnnotationPrefix + "cpuset_trick_exempt_tasks"
	SpecCPUBvtWarpNs              = AnnotationPrefix + "cpu_bvt_warp_ns"
	SpecCpuacctSchedLatSwitch     = AnnotationPrefix + "cpuacct_sched_lat_switch"
	SpecMemoryWmarkRatio          = AnnotationPrefix + "memory_wmark_ratio"
	SpecMemoryExtra               = AnnotationPrefix + "memory_extra"
	SpecMemoryForceEmptyCtl       = AnnotationPrefix + "memory_force_empty_ctl"
	SpecMemoryPriority            = AnnotationPrefix + "memory_priority"
	SpecMemoryUsePriorityOOM      = AnnotationPrefix + "memory_use_priority_oom"
	SpecMemoryOOMKillAll          = AnnotationPrefix + "memory_oom_kill_all"
	SpecMemoryDroppable           = AnnotationPrefix + "memory_droppable"
	SpecIntelRdtL3Cbm             = AnnotationPrefix + "intel_rdt_l3_cbm"
	SpecIntelRdtGroup             = AnnotationPrefix + "intel_rdt_group"
	SpecIntelRdtMba               = AnnotationPrefix + "intel_rdt_mba"
	SpecBlkioFileLevelSwitch      = AnnotationPrefix + "blkio_file_level_switch"
	SpecBlkioBufferWriteBps       = AnnotationPrefix + "blkio_buffer_write_bps"
	SpecBlkioMetaWriteTps         = AnnotationPrefix + "blkio_meta_write_tps"
	SpecBlkioFileThrottlePath     = AnnotationPrefix + "blkio_file_throttle_path"
	SpecBlkioBufferWriteSwitch    = AnnotationPrefix + "blkio_buffer_write_switch"
	SpecBlkioDeviceBufferWriteBps = AnnotationPrefix + "blkio_device_buffer_write_bps"
	SpecBlkioDeviceIdleTime       = AnnotationPrefix + "blkio_device_idle_time"
	SpecBlkioDeviceLatencyTarget  = AnnotationPrefix + "blkio_device_latency_target"
	SpecBlkioDeviceReadLowBps     = AnnotationPrefix + "blkio_device_read_low_bps"
	SpecBlkioDeviceReadLowIOps    = AnnotationPrefix + "blkio_device_read_low_iops"
	SpecBlkioDeviceWriteLowBps    = AnnotationPrefix + "blkio_device_write_low_bps"
	SpecBlkioDeviceWriteLowIOps   = AnnotationPrefix + "blkio_device_write_low_iops"
	SpecNetCgroupRate             = AnnotationPrefix + "net_cgroup_rate"
	SpecNetCgroupCeil             = AnnotationPrefix + "net_cgroup_ceil"
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
