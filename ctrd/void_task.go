package ctrd

import (
	"context"
	"fmt"
	"syscall"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/api/types"
	"github.com/containerd/containerd/cio"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type voidTask struct {
	id string
}

func newVoidTask(id string) (containerd.Task, error) {
	return &voidTask{
		id: id,
	}, nil
}

// ID of the process
func (v *voidTask) ID() string {
	return v.id
}

// Pid is the system specific process id
func (v *voidTask) Pid() uint32 {
	return 0
}

// Start starts the process executing the user's defined binary
func (v *voidTask) Start(context.Context) error {
	return fmt.Errorf("can not start the process in void task")
}

// Delete removes the process and any resources allocated returning the exit status
func (v *voidTask) Delete(context.Context, ...containerd.ProcessDeleteOpts) (*containerd.ExitStatus, error) {
	return nil, fmt.Errorf("can not delete process inside void task")
}

// Kill sends the provided signal to the process
func (v *voidTask) Kill(context.Context, syscall.Signal, ...containerd.KillOpts) error {
	return nil
}

// Wait asynchronously waits for the process to exit, and sends the exit code to the returned channel
func (v *voidTask) Wait(context.Context) (<-chan containerd.ExitStatus, error) {
	return nil, fmt.Errorf("can not wait the void task to exit")
}

// CloseIO allows various pipes to be closed on the process
func (v *voidTask) CloseIO(context.Context, ...containerd.IOCloserOpts) error {
	return nil
}

// Resize changes the width and heigh of the process's terminal
func (v *voidTask) Resize(ctx context.Context, w, h uint32) error {
	return nil
}

// IO returns the io set for the process
func (v *voidTask) IO() cio.IO {
	return nil
}

// Status returns the executing status of the process
func (v *voidTask) Status(context.Context) (containerd.Status, error) {
	return containerd.Status{}, fmt.Errorf("can not Status for void Task")
}

// Pause suspends the execution of the task
func (v *voidTask) Pause(context.Context) error {
	return fmt.Errorf("can not Pause a void Task")
}

// Resume the execution of the task
func (v *voidTask) Resume(context.Context) error {
	return nil
}

// Exec creates a new process inside the task
func (v *voidTask) Exec(context.Context, string, *specs.Process, cio.Creator) (containerd.Process, error) {
	return nil, fmt.Errorf("can not create a new process inside a void task")
}

// Pids returns a list of system specific process ids inside the task
func (v *voidTask) Pids(context.Context) ([]containerd.ProcessInfo, error) {
	return []containerd.ProcessInfo{}, nil
}

// Checkpoint serializes the runtime and memory information of a task into an
// OCI Index that can be push and pulled from a remote resource.
//
// Additional software like CRIU maybe required to checkpoint and restore tasks
func (v *voidTask) Checkpoint(context.Context, ...containerd.CheckpointTaskOpts) (containerd.Image, error) {
	return nil, fmt.Errorf("can not Checkpoint a void task")
}

// Update modifies executing tasks with updated settings
func (v *voidTask) Update(context.Context, ...containerd.UpdateTaskOpts) error {
	return nil
}

// LoadProcess loads a previously created exec'd process
func (v *voidTask) LoadProcess(context.Context, string, cio.Attach) (containerd.Process, error) {
	return nil, nil
}

// Metrics returns task metrics for runtime specific metrics
//
// The metric types are generic to containerd and change depending on the runtime
// For the built in Linux runtime, github.com/containerd/cgroups.Metrics
// are returned in protobuf format
func (v *voidTask) Metrics(context.Context) (*types.Metric, error) {
	return nil, nil
}
