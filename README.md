# mini-container

学习容器技术时写的玩具，目前实现了：

- 用 namespace 实现资源隔离
- 用 cgroups 实现资源限制
- 挂载新的 rootfs

还没学：

- 文件系统的高级玩法
- AppArmor 等安全机制

如果你刚好有测试环境并且想玩一玩：

```bash
# Ubuntu 20.04LTS amd64
# Kernel 5.4.0-105-generic
# go1.17

go build && sudo ./mini-container
```
