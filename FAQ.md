# FAQ

## Why does my console not support line wrapping or color?

`vm` connects to domains using a virtual serial TTY device. From the guest's
perspective, each serial TTY that's attached to the guest gets a corresponding
`agetty` process to manage it. That process determines the value of `TERM`
variable that gets passed into the environment of the resulting shell session.

Typically, the init system of a Linux distribution sets up the agetty process
spawning mechanism for serial TTYs. For systemd, those TTYs are managed by
instantiated `serial-getty@.service` units. systemd sets the value of `TERM`
for its serial TTYs to `vt220` for [compatibility reasons](https://github.com/systemd/systemd/issues/3342#issuecomment-221821337).

To override the `TERM` value for serial TTYs, create the file `/etc/systemd/system/serial-getty@.service.d/override.conf`
with the contents:

```
[Service]
Environment=TERM=xterm
```
