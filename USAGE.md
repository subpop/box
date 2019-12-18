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

**--disk**="": Attach `FILE` to the domain as a secondary disk

**--name**="": Assign `NAME` to the domain

**--no-snapshot**: Disable taking an initial snapshot upon creation

**--transient**: Create a non-persistent domain

## list

List defined domains

**--all**: Include inactive domains

**--inactive**: List only inactive domains

## destroy

Destroy a domain

**--force**: Immediately destroy the domain, without prompting

## up

Start a domain

**--connect**: Immediately connect to the started domain

## down

Stop a domain

**--force**: Immediately stop the domain, without prompting

**--graceful**: Power off the domain gracefully

## restart

Restart a domain

**--force**: Immediately restart the domain, without prompting

**--graceful**: Restart the domain gracefully

## connect

Connect to a running domain

**--identity**="": Attempt SSH authentication using `IDENTITY`

**--mode**="": Connection mode: serial, console, or ssh (default: serial)

**--user**="": User to connect as over SSH (default: root)

## inspect

Show details about a domain

**--format**="": Specify output format

## image

Manage backing disk images

### list

List available backing disk images

### get

Retrieve a new backing disk image

**--quiet**: No progress output

**--rename**="": Rename backing disk image to `NAME`

### remove

Remove a backing disk image

**--force**: Force removal of a backing disk image without prompting

## template

Manage backing disk templates from libguestfs

### list

List templates available for import

**--sort**="": Sort list by `VALUE` (default: name)

### sync

Refresh available templates from build service

### info

Print details about a template

**--arch**="": Specify alternate architecture (default: x86_64)

### get

Retrieve and prepare a template from build service

**--arch**="": Specify alternative architecture (default: x86_64)

**--quiet**: No progress output

## snapshot

Manage domain snapshots

### list

List snapshots for a domain

### create

Take a new snapshot for a domain

**--name**="": Create a snapshot with `NAME`

### remove

Remove a snapshot for a domain

**--snapshot**="": Remove snapshot named `NAME`

### revert

Revert a domain to snapshot

**--snapshot**="": Revert to `SNAPSHOT`

## help, h

Shows a list of commands or help for one command

