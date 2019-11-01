# About #

`vm` is a command line utility that provides a high-level interface to create
and manage virtual machines through libvirt.

# Installation #

```bash
go get -u github.com/subpop/vm/cmd/vm
```

# Usage #

Download a base image:

```bash
vm image get https://dl.fedoraproject.org/pub/fedora/linux/releases/31/Cloud/x86_64/images/Fedora-Cloud-Base-31-1.9.x86_64.qcow2
```

Download and convert a Vagrant ".box":

```bash
vm image get https://dl.fedoraproject.org/pub/fedora/linux/releases/31/Cloud/x86_64/images/Fedora-Cloud-Base-Vagrant-31-1.9.x86_64.vagrant-libvirt.box
```

Create a domain backed by that image:

```bash
vm create Fedora-Cloud-Base-31-1.9.x86_64 --name my-f31
```

List active domains:

```bash
vm list
```

Start a created domain:

```bash
vm up my-f31
```

Connect to an existing domain over SSH:

```bash
vm connect -m ssh -u vagrant my-f31
```

Connect to an existing domain over VirtIO PTY:

```bash
vm connect -m console my-f31
```
