### viu

This is an example in Rust from https://github.com/wapm-packages/viu using `tonistiigi/xx:rust` base image.

Easiest to follow the example by running `make shell` followed by `cd hello`.

```
docker buildx create --use
docker buildx build --platform=linux/amd64,wasi/wasm -t tonistiigi/viu --push .
docker buildx build --platform=linux/amd64,wasi/wasm -t tonistiigi/viu:docker --target=docker --push .
docker buildx imagetools inspect tonistiigi/viu
docker run tonistiigi/viu:docker
docker wasm run tonistiigi/viu:docker
docker wasm run --runtime=wasmer tonistiigi/viu:docker
docker wasm run --runtime=wasmer tonistiigi/viu /mindblown.gif
```
