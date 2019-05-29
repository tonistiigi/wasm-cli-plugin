### gohello

This is an example in Go. WASI support in Go is in a very early development phase. Do not expect anything but "hello-world" to work. It also seems that go binaries do not currently work in wasmer runtime.

```
docker buildx create --use
docker buildx build --platform=linux/amd64,wasi/wasm -t tonistiigi/hello:go --push .
docker buildx imagetools inspect tonistiigi/hello:go
docker run tonistiigi/hello:go
docker wasm run -e FOO=bar tonistiigi/hello:go
```