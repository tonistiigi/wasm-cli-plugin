```
docker buildx create --use
docker buildx build --platform=linux/amd64,wasi/wasm -t myuser/viu --push .
docker buildx build --platform=linux/amd64,wasi/wasm -t myuser/viu:docker --target=docker --push .
docker buildx imagetools inspect myuser/viu
docker run myuser/viu:docker
docker wasm run myuser/viu:docker
docker wasm run --runtime=wasmer myuser/viu:docker
docker wasm run --runtime=wasmer myuser/viu /success.gif
```