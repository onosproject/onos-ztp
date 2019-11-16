# How to install and run onos-ztp?

The current implementation of the ZTP subsystem uses local file system storage for the role configuration; this is only temporary.

To start the server simply run:
```bash
> go run github.com/onosproject/onos-ztp/cmd/onos-ztp
```

You may then use the ONOS CLI `ztp` commands from the consolidated ONOS CLI program, e.g.

```bash
> onos ztp get roles
> onos ztp add role test/samplejson/leaf.json
> onos ztp add role test/samplejson/spine.json
> onos ztp get role leaf
> onos ztp remove role spine
```

_More documentation to be added._

