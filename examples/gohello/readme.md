```
docker buildx create --use
docker buildx build --platform=linux/amd64,wasi/wasm -t tonistiigi/hello:go --push .
docker buildx imagetools inspect tonistiigi/hello:go
docker run tonistiigi/hello:go
docker wasm run -e FOO=bar tonistiigi/hello:go
```