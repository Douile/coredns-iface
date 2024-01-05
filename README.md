# iface

## Name

*iface* - Returns address of network interface associated with provided name

## Description

The plugin looks up a network interface using the provided name and the configured address for
that interface if it exists.

Note: In order to add a new plugin, an additional step of `make gen` is needed. Therefore,
to build the coredns with demo plugin the following should be used:
```
docker run -it --rm -v $PWD:/v -w /v golang:1.16 sh -c 'make gen && make'
```

## Syntax

~~~ txt
iface
~~~

## Also See

See the [manual](https://coredns.io/manual).
