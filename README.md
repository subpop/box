# About #

`box` is a command line utility that provides a high-level interface to create
and manage virtual machines through libvirt.

# Usage #

Download a base image:

```bash
box image get -n fedora-30
```

Create a box from that image:

```bash
box create -i fedora-30
```

List available boxes:

```bash
box list
```

Start a created box:

```bash
box up -n awaited-sawfly
```

Connect to an existing box over SSH:

```bash
box connect -m ssh -n awaited-sawfly
```