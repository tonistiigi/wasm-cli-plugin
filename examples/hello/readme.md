```
docker buildx create --use
docker buildx build --platform=linux/amd64,linux/arm64,linux/arm,wasi/wasm -t myuser/hello --push .
docker buildx imagetools inspect myuser/hello
docker run myuser/hello
docker wasm run myuser/hello
```