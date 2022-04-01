//go:build linux
// +build linux

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if len(os.Args) > 1 {
		initContainer()
		return
	}

	c := exec.Command("/proc/self/exe", "init")
	// 创建namespace
	c.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET,
		// Credential: &syscall.Credential{Uid: uint32(1), Gid: uint32(1)},
	}
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	err := c.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func initContainer() {
	setRootfs()
	setCgroups()

	err := syscall.Mount("proc", "/proc", "proc", syscall.MS_NODEV|syscall.MS_NOEXEC|syscall.MS_NOSUID, "")
	if err != nil {
		log.Fatalln(err)
	}
	err = syscall.Sethostname([]byte("container-demo"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("==== welcome to a new space ====")
	// 通过execve系统调用让bash代替当前程序成为pid=1进程
	err = syscall.Exec("/bin/bash", nil, os.Environ())
	if err != nil {
		log.Fatalln(err)
	}
}
