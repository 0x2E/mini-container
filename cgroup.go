//go:build linux
// +build linux

package main

import (
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

const (
	cgroutRoot = "/sys/fs/cgroup"
	groupName  = "container_demo"

	// 最大占用50MB内存
	limitInBytes = "50m"
	// 每100ms可用10ms CPU
	cfsQuotaUs  = "10000"
	cfsPeriodUs = "100000"
)

func setCgroups() {
	pid := os.Getpid()
	// mem
	// memCg := path.Join(cgroutRoot, "memory", groupName)
	// os.Mkdir(memCg, 0755)
	// os.WriteFile(path.Join(memCg, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	// os.WriteFile(path.Join(memCg, "memory.limit_in_bytes"), []byte("50M"), 0644)
	memCg := newCgroup("memory", pid)
	memCg.set("memory.limit_in_bytes", limitInBytes)

	// cpu
	cpuCg := newCgroup("cpu", pid)
	cpuCg.set("cpu.cfs_period_us", cfsPeriodUs)
	cpuCg.set("cpu.cfs_quota_us", cfsQuotaUs)
}

type cgroup struct {
	path     string
	task     []byte
	original map[string][]byte
}

func newCgroup(subsystem string, pid int) *cgroup {
	return &cgroup{
		path:     path.Join(cgroutRoot, subsystem, groupName),
		task:     []byte(strconv.Itoa(pid)),
		original: make(map[string][]byte),
	}
}

func (c *cgroup) set(key, data string) {
	err := os.MkdirAll(c.path, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(path.Join(c.path, "tasks"), c.task, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.OpenFile(path.Join(c.path, key), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln(err)
	}
	defer f.Close()
	original, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = f.WriteString(data)
	if err != nil {
		log.Fatalln(err)
	}
	c.original[key] = []byte(original)
}

// func (c *cgroup) unset(key string) {
// 	original, ok := c.original[key]
// 	if !ok {
// 		log.Println("cgroup key not exists")
// 		return
// 	}
// 	err := os.WriteFile(path.Join(c.path, key), original, 0644)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	delete(c.original, key)
// }

// func (c *cgroup) destory() {
// 	err := syscall.Rmdir(c.path)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
