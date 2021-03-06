//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package runtime

import (
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/lastbackend/lastbackend/pkg/util/system"
	"github.com/lastbackend/registry/pkg/builder/envs"
	"github.com/lastbackend/registry/pkg/distribution/types"
	"github.com/shirou/gopsutil/mem"
)

func BuilderInfo() types.BuilderInfo {

	hostname, err := os.Hostname()
	if err != nil {
		_ = fmt.Errorf("get hostname err: %v", err)
	}

	info := types.BuilderInfo{}
	osInfo := system.GetOsInfo()
	info.Hostname = hostname
	info.OSType = osInfo.GoOS
	info.OSName = fmt.Sprintf("%s %s", osInfo.OS, osInfo.Core)
	info.Architecture = osInfo.Platform
	return info
}

func BuilderStatus() types.BuilderStatus {
	state := types.BuilderStatus{}
	state.Capacity = BuilderCapacity()
	state.Usage = BuilderUsage()
	return state
}

func BuilderCapacity() types.BuilderResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		_ = fmt.Errorf("get memory err: %v", err)
	}

	var stat syscall.Statfs_t

	wd, err := os.Getwd()

	err = syscall.Statfs(wd, &stat)
	if err != nil {
		_ = fmt.Errorf("get stats err: %v", err)
	}

	// Available blocks * size per block = available space in bytes
	storage := stat.Bfree * uint64(stat.Bsize)

	m := vmStat.Total / 1024 / 1024

	return types.BuilderResources{
		Storage: int64(storage) / 1024 / 1024,
		Workers: int(m / types.DEFAULT_MIN_WORKER_MEMORY),
		RAM:     int64(m),
		CPU:     int64(runtime.NumCPU()),
	}
}

func BuilderUsage() types.BuilderResources {

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		_ = fmt.Errorf("get memory err: %v", err)
	}

	m := vmStat.Free / 1024 / 1024

	return types.BuilderResources{
		Workers: envs.Get().GetBuilder().ActiveWorkers(),
		RAM:     int64(m),
		CPU:     0,
	}
}
