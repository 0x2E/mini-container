//go:build linux
// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func setRootfs() {
	// systemd会在启动时把所有挂载点的传播类型改为share，所以这里改为private防止影响其他namespace
	// https://man7.org/linux/man-pages/man7/mount_namespaces.7.html#NOTES
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.MkdirAll("tmp/rootfs/rootfs_old", 0755)
	if err != nil {
		log.Fatalln(err)
	}
	// https://cdimage.ubuntu.com/ubuntu-base/releases/20.04/release/ubuntu-base-20.04.4-base-amd64.tar.gz
	// 测试使用alpine做rootfs时会crash
	err = exec.Command("tar", "-zxf", "ubuntu-base-20.04.4-base-amd64.tar.gz", "-C", "tmp/rootfs").Run()
	if err != nil {
		log.Fatalln(err)
	}
	// 挂载自己以创建新的挂载点，满足PivotRoot的条件
	err = syscall.Mount("tmp/rootfs", "tmp/rootfs", "", syscall.MS_BIND, "")
	if err != nil {
		log.Fatalln(err)
	}
	err = syscall.PivotRoot("tmp/rootfs", "tmp/rootfs/rootfs_old")
	if err != nil {
		log.Fatalln(err)
	}
	// 现在已经切换到了新的rootfs
	err = syscall.Chdir("/")
	if err != nil {
		log.Fatalln(err)
	}
	err = syscall.Unmount("/rootfs_old", syscall.MNT_DETACH)
	if err != nil {
		log.Fatalln(err)
	}
}
