# onos-ztp
[![Build Status](https://travis-ci.org/onosproject/onos-ztp.svg?branch=master)](https://travis-ci.org/onosproject/onos-ztp)
[![Go Report Card](https://goreportcard.com/badge/github.com/onosproject/onos-ztp)](https://goreportcard.com/report/github.com/onosproject/onos-ztp)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/gojp/goreportcard/blob/master/LICENSE)
[![Coverage Status](https://img.shields.io/coveralls/github/onosproject/onos-ztp/badge.svg)](https://coveralls.io/github/onosproject/onos-ztp?branch=master)
[![GoDoc](https://godoc.org/github.com/onosproject/onos-ztp?status.svg)](https://godoc.org/github.com/onosproject/onos-ztp)

Zero-Touch Provisioning subsystem built using the µONOS architecture.

## Design Objectives
Setting up and managing network infrastructure devices often requires elaborate procedures 
to be followed in order to maintain network integrity. The goal of this subsystem is to 
ease this burden for the operators and to make the lifecycle management of network 
infrastructure devices simpler, faster and more predictable.

This is accomplished by allowing the operators to predefine various classes or roles 
for the devices in the network and manage their configurations and pipeline definitions
on per-class basis, thus increasing consistency and reducing toil.

This subsystem allows the operators to manage the role class configurations and 
pipeline definitions and in turn apply to the network devices through the `onos-config` and `onos-control` 
subsystems, respectively.

_More documentation to be added._

## Running onoz-ztp

The current implementation of the ZTP subsystem uses local file system storage for the role configuration; this is only temporary.

To start the server simply run:
```
> go run github.com/onosproject/onos-ztp/cmd/onos-ztp
```

You may then use the ONOS CLI `ztp` commands from the consolidated ONOS CLI program, e.g.

```
> onos ztp get roles
> onos ztp add role test/samplejson/leaf.json
> onos ztp add role test/samplejson/spine.json
> onos ztp get role leaf
> onos ztp remove role spine
```
