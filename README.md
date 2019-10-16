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
vm image get -n fedora-30
```

Create a VM from that image:

```bash
vm create -i fedora-30
```

List available VMs:

```bash
vm list
```

Start a created VM:

```bash
vm up -n awaited-sawfly
```

Connect to an existing VM over SSH:

```bash
vm connect -m ssh -n awaited-sawfly
```

Connect to an existing VM over console TTY:

```bash
vm connect -m console -n awaited-sawfly
```
