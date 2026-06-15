# Star Conflict Masterserver Reimplementation

## Instructions

Requirements:
- [podman](https://podman.io/docs/installation)
- [golang](https://go.dev/doc/install)

### Linux Instructions

### Windows Instructions

[Documentdb](https://documentdb.io/) only does releases for Linux, so a Windows Host with Podman Desktop and enabled Virtualization is required.[^1]

## Components

### Loadbalancer

The initial Service the game connects to.
Provides some initial configuration and points the game to the assigned shard and chat server.

### Shard
 
The core masterserver component, handles everything that happens between login and logout.

[^1]: Subject to change at a later date
