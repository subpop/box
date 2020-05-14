% vm(8) `vm` is a program to manage and interact with virtual machines.


# NAME

vm - control virtual machines

# SYNOPSIS

vm

```
[--generate-fish-completion]
[--generate-man-page]
[--generate-markdown]
[--help|-h]
[--version|-v]
```

**Usage**:

```
vm [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--generate-fish-completion**: 

**--generate-man-page**: 

**--generate-markdown**: 

**--help, -h**: show help

**--version, -v**: print the version


# COMMANDS

## create

Creates a new domain from the specified image

**--detach**: Detach from the newly created domain

**--disk, -d**="": Attach `FILE` to the domain as a secondary disk

**--memory, -m**="": Create a domain with `MEM` RAM (default: 256 MB)

**--name, -n**="": Assign `NAME` to the domain

**--network, -N**="": Use bridged network device `BRIDGE`

**--no-snapshot**: Disable taking an initial snapshot upon creation

**--transient, -t**: Create a non-persistent domain

**--uefi**: Use UEFI boot loader

**--video, -v**="": Use video device `TYPE`

## list

List defined domains

**--all**: Include inactive domains

**--inactive**: List only inactive domains

## destroy

Destroy a domain

**--force, -f**: Immediately destroy the domain, without prompting

## up

Start a domain

**--connect, -c**: Immediately connect to the started domain

## down

Stop a domain

**--force, -f**: Immediately stop the domain, without prompting

**--graceful, -g**: Power off the domain gracefully

## restart

Restart a domain

**--force, -f**: Immediately restart the domain, without prompting

**--graceful, -g**: Restart the domain gracefully

## connect

Connect to a running domain

**--identity, -i**="": Attempt SSH authentication using `IDENTITY`

**--mode, -m**="": Connection mode: serial, console, or ssh (default: serial)

**--user, -u**="": User to connect as over SSH (default: root)

## info

Show details about a domain

## dump

Show XML description of a domain

## image

Manage backing disk images

### list

List available backing disk images

### get

Retrieve a new backing disk image

**--quiet, -q**: No progress output

**--rename, -r**="": Rename backing disk image to `NAME`

### remove

Remove a backing disk image

**--force, -f**: Force removal of a backing disk image without prompting

## template

Manage backing disk templates from libguestfs

### list

List templates available for import

**--sort, -s**="": Sort list by `VALUE` (default: name)

### sync

Refresh available templates from build service

### info

Print details about a template

**--arch, -a**="": Specify alternate architecture (default: x86_64)

### get

Retrieve and prepare a template from build service

**--arch, -a**="": Specify alternative architecture (default: x86_64)

**--quiet, -q**: No progress output

## snapshot

Manage domain snapshots

### list

List snapshots for a domain

### create

Take a new snapshot for a domain

**--name, -n**="": Create a snapshot with `NAME`

### remove

Remove a snapshot for a domain

**--snapshot, -s**="": Remove snapshot named `NAME`

### revert

Revert a domain to snapshot

**--snapshot, -s**="": Revert to `SNAPSHOT`

## capabilities

Get details on hypervisor capabilities

**--format, -f**="": Specify output format (default: xml)

## domain-capabilities

Get details on domain capabilities

**--format, -f**="": Specify output format (default: xml)

## help, h

Shows a list of commands or help for one command

