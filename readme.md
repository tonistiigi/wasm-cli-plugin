## wasm-cli-plugin

This is a Docker CLI plugin that allows running `wasi/wasm` applications built from Dockerfiles and distributed with Docker registries on any platform as an unprivileged user.

![demo](https://imgur.com/download/pvBRyAi)

For more context, watch the "Docker + WebAssembly" video from the [WasmSF](https://www.meetup.com/wasmsf/) meetup https://www.youtube.com/watch?v=rZEQNH27y-k .

### Getting started

Running `make shell` puts you in a development environment where you have access to `docker`, `buildx` and the example apps in this repository.

#### Available commands

The UI is similar to docker but simplified a lot. You can pull an image with `docker wasm pull <image>` and run with `docker wasm run <image> <cmd>`. 

```
# docker wasm

Usage:	docker wasm [OPTIONS] COMMAND

Run wasm containers

Options:
      --data-root string   Root directory of persistent state (default "/root/.docker/wasm")
  -D, --debug              Enable debug logs

Commands:
  ls          List wasm images
  pull        Pull a wasm image
  rm          Remove a wasm image
  run         Run a wasm image
```

### Building

```
make binaries
```

Builds `docker-wasm` binary under `./bin` folder. 

### Installing

Build `docker-wasm` or pull from [releases](https://github.com/tonistiigi/wasm-cli-plugin/releases) place it on `PATH`. Either [`wasmtime`](https://github.com/CraneStation/wasmtime) or [`wasmer`](https://github.com/wasmerio/wasmer/releases) runtime needs to be on `PATH` to run wasi applications. On Linux and OSX `wasmtime` runtime binary is also built with `docker-wasm`. 

`docker-wasm` binary can be used directly or since `v19.03` it can be used as a Docker CLI plugin, adding a `docker wasm` command to Docker CLI. To install the plugin copy the built binary to `~/.docker/cli-plugins/docker-wasm`.


### Building multi-platform images with wasm support

Wasm applications are built from Dockerfiles by setting `--platform=wasi/wasm` flag on build. `--platform` is available (without experimental) since Docker 19.03 (only with BuildKit). Docker [buildx](https://github.com/docker/buildx) plugin is another way to build and allows specifying multiple platforms together into a single multi-platform image and works with previous Docker engine versions.

[Automatic platform args](https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope) can be used in Dockerfiles to efficiently cross-compile  to any architecture with multi-stage builds. https://github.com/tonistiigi/xx repository contains a set of base images that automatically integrate with platform args, requiring no configuration from user to automatically pick the correct toolchain for a specific `--platform` value.

The `./examples` directory contains example applications. Every application has its  readme with instructions.

#### [hello](https://github.com/tonistiigi/wasm-cli-plugin/tree/master/examples/hello) - hello world application in c, demonstrating POSIX capabilities

#### [viu](https://github.com/tonistiigi/wasm-cli-plugin/tree/master/examples/viu) - image preview application in rust

#### [gohello](https://github.com/tonistiigi/wasm-cli-plugin/tree/master/examples/gohello) - hello world in go

### Alternatives: containerd-shim

Multi-platform images with `wasi/wasm` support can be also run in `containerd` using [containerd wasm shim](https://github.com/dmcgowan/containerd-wasm).

Example:
```
ctr image pull --platform wasi/wasm docker.io/tonistiigi/hello:latest
ctr run --rm --runtime io.containerd.wasm.v1 --platform=wasi/wasm docker.io/tonistiigi/hello:latest hello
```

### Alternatives: Docker image

`docker wasm` can also be tested with a docker image `tonistiigi/docker-wasm`.

Example:
```
docker run --rm tonistiigi/docker-wasm run tonistiigi/viu /success.gif
```
