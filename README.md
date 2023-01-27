# DREAMEMO

> You Can (Not) Escape

![]()

## Architecture

![dreamemo-arch](./image/dreamemo-arch.png)

## Quick Start

- RPC data protocol
  - thrift
  - protobuf
- Cache elimination algorithm (use interface to abstract)
  - LRU
  - LFU
- Distributed algorithm (use interface to abstract)
  - Consistent Hash
  - Raft
- Supported Data Source (use interface to abstract)
  - Redis

User can choose to use these features and the numbers of nodes

TODO: Add an interactive command line to allow the user to select a configuration.

DREAMEMO is a subproject of the [BINARY WEB ECOLOGY](https://github.com/B1NARY-GR0UP)