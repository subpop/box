# About #

`vm` is a command line utility that provides a high-level interface to create
and manage virtual machines through libvirt.

It supports creation and deletion of domains and snapshots, connecting to serial
and SSH terminals, downloading and converting templates from vagrantup.com and
builder.libguestfs.org, and importing locally downloading images.

# Installation #

```bash
go get -u github.com/subpop/vm/cmd/vm
```

# Distribution #

`vm` includes a `Makefile` to aid distributions in packaging. The default target
will build `vm` along with shell completion data and a man page. The `Makefile`
includes an `install` target to install the binary and data into distribution-appropriate
locations. To override the installation directory (commonly referred to as the
`DESTDIR`), set the `DESTDIR` variable when running the `install` target:

```bash
[link@localhost vm]$ make
go build ./cmd/vm
go run ./cmd/vm --generate-fish-completion > vm.fish
go run ./cmd/vm --generate-bash-completion > vm.bash
go run ./cmd/vm --generate-man-page > vm.1
gzip -k vm.1
[link@localhost vm]$ make DESTDIR=_inst install
install -D -m755 -t _inst//usr/local/bin vm
install -D -m644 -t _inst//usr/local/share/man/man1 vm.1.gz
install -D -m644 -t _inst//usr/local/share/fish/completions vm.fish
install -d _inst//usr/local/share/bash-completion/completions
install -m644 -T vm.bash _inst//usr/local/share/bash-completion/completions/vm
[link@localhost vm]$
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

Create a domain without defining it:

```bash
vm create Fedora-Cloud-Base-31-1.9.x86_64 --transient
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
vm connect -m ssh -u vagrant my-f31 -i ~/.ssh/cloud_user_rsa
```

Connect to an existing domain over VirtIO PTY:

```bash
vm connect -m console my-f31
```

Take a snapshot:

```bash
vm snapshot create my-f31 --name fresh_install
```

Revert to snapshot:

```bash
vm snapshot revert my-f31 --snapshot fresh_install
```
