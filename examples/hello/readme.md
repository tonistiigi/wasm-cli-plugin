### hello

This is an example in C demonstrating POSIX capabilities.

Easiest to follow the example by running `make shell` followed by `cd hello`.

[Buildx](https://github.com/docker/buildx) is needed to create a multi-platform image that has both wasm subimages and linux container images.

Switch buildx to `docker-container` driver if it isn't already:

```
docker buildx create --use
```

Build application for amd64, arm64 and wasi, pushing the result to the registry.

```
> docker buildx build --platform=linux/amd64,linux/arm64,linux/arm,wasi/wasm -t tonistiigi/hello --push .

/work/hello # docker buildx build --platform=linux/amd64,linux/arm64,linux/arm,wasi/wasm -t tonistiigi/hello --push .
[+] Building 4.2s (22/22) FINISHED
 => [internal] load .dockerignore                                                                                        0.0s
 => => transferring context: 2B                                                                                          0.0s
 => [internal] load build definition from Dockerfile                                                                     0.0s
 => => transferring dockerfile: 32B                                                                                      0.0s
 => [linux/amd64 internal] load metadata for docker.io/tonistiigi/xx:llvm                                                0.4s
 => [linux/amd64 build 1/4] FROM docker.io/tonistiigi/xx:llvm@sha256:9245fa10d2cd2bdabfa2f0aacd54a7a5bf25656a66b84d80c1  0.0s
 => => resolve docker.io/tonistiigi/xx:llvm@sha256:9245fa10d2cd2bdabfa2f0aacd54a7a5bf25656a66b84d80c1e969b9a2018c90      0.0s
 => [internal] load build context                                                                                        0.0s
 => => transferring context: 29B                                                                                         0.0s
 => CACHED [linux/amd64 build 2/4] WORKDIR /src                                                                          0.0s
 => CACHED [linux/amd64 build 3/4] COPY hello.c .                                                                        0.0s
 => CACHED [linux/amd64 build 4/4] RUN clang -static -o /hello hello.c                                                   0.0s
 => CACHED [linux/arm/v7 stage-1 1/1] COPY --from=build /hello .                                                         0.0s
 => CACHED [linux/amd64 build 2/4] WORKDIR /src                                                                          0.0s
 => CACHED [linux/amd64 build 3/4] COPY hello.c .                                                                        0.0s
 => CACHED [linux/amd64 build 4/4] RUN clang -static -o /hello hello.c                                                   0.0s
 => CACHED [linux/arm64 stage-1 1/1] COPY --from=build /hello .                                                          0.0s
 => CACHED [linux/amd64 build 2/4] WORKDIR /src                                                                          0.0s
 => CACHED [linux/amd64 build 3/4] COPY hello.c .                                                                        0.0s
 => CACHED [linux/amd64 build 4/4] RUN clang -static -o /hello hello.c                                                   0.0s
 => CACHED [linux/amd64 stage-1 1/1] COPY --from=build /hello .                                                          0.0s
 => CACHED [linux/amd64 build 2/4] WORKDIR /src                                                                          0.0s
 => CACHED [linux/amd64 build 3/4] COPY hello.c .                                                                        0.0s
 => CACHED [linux/amd64 build 4/4] RUN clang -static -o /hello hello.c                                                   0.0s
 => CACHED [wasi/wasm stage-1 1/1] COPY --from=build /hello .                                                            0.0s
 => exporting to image                                                                                                   3.7s
 => => exporting layers                                                                                                  0.0s
 => => exporting manifest sha256:20385c23e201b0cf330dcf2b1d637947b08400f737d464d2563d471a52a080e2                        0.0s
 => => exporting config sha256:d68e98c59f84d741a67a3aafe32f690b24a58882c68b2d5670727ccb77b2b128                          0.0s
 => => exporting manifest sha256:cda320d2dc2ff1218f439288cfb239cd1b452e51409828c3acfe0c47a27a652e                        0.0s
 => => exporting config sha256:4a19321cc5b3ae85c826fc8fe4280302129085b7e5f9b2851b370b3bc700fea5                          0.0s
 => => exporting manifest sha256:8ac9809fa588af234b904ba3dcaf6472af8bc49f45836c73424ddec90030ac59                        0.0s
 => => exporting config sha256:7a8580e6977c1b4f16303467b71ca7735a5f2b78ee708330be45c4bb64473d0a                          0.0s
 => => exporting manifest sha256:d020d70549e29936cb929a3c8c0fa4f1d91f8b349d0f8512d8b37a449884d713                        0.0s
 => => exporting config sha256:cf8192ff388df08ea6446cc8e0046d5d33653108ada300a7436ce8dc8a0e4338                          0.0s
 => => exporting manifest list sha256:15163c37144fc9dc92e8e7417eca90055285aa716786ad70a471dbf40f51fa75                   0.0s
 => => pushing layers                                                                                                    0.9s
 => => pushing manifest for docker.io/tonistiigi/hello:latest                                                            2.7s
```

Inspect the multi-platform image in registry. Notice the three platforms we built for.

```
> docker buildx imagetools inspect tonistiigi/hello

Name:      docker.io/tonistiigi/hello:latest
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:15163c37144fc9dc92e8e7417eca90055285aa716786ad70a471dbf40f51fa75

Manifests:
  Name:      docker.io/tonistiigi/hello:latest@sha256:20385c23e201b0cf330dcf2b1d637947b08400f737d464d2563d471a52a080e2
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/amd64

  Name:      docker.io/tonistiigi/hello:latest@sha256:cda320d2dc2ff1218f439288cfb239cd1b452e51409828c3acfe0c47a27a652e
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm64

  Name:      docker.io/tonistiigi/hello:latest@sha256:8ac9809fa588af234b904ba3dcaf6472af8bc49f45836c73424ddec90030ac59
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm/v7

  Name:      docker.io/tonistiigi/hello:latest@sha256:d020d70549e29936cb929a3c8c0fa4f1d91f8b349d0f8512d8b37a449884d713
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  wasi/wasm
```

Run as a regular linux container:

```
docker run tonistiigi/hello

Hello world, I am linux/amd64!
contents of / :
  4	.
  4	..
  8	hello
  8	.dockerenv
  4	etc
  4	dev
  4	sys
  4	proc

writing to /foo
contents of / after write :
  4	.
  4	..
  8	hello
  8	foo
  8	.dockerenv
  4	etc
  4	dev
  4	sys
  4	proc
```

Run with webassembly:

```
docker wasm run tonistiigi/hello

docker wasm run tonistiigi/hello
INFO[0000] pulling sha256:15163c37144fc9dc92e8e7417eca90055285aa716786ad70a471dbf40f51fa75
INFO[0001] pulling sha256:d020d70549e29936cb929a3c8c0fa4f1d91f8b349d0f8512d8b37a449884d713
INFO[0001] pulling sha256:bd0c74e37268f0ad2d91d57c2547f9c007a263e9aff75514772a5d5ca8305577
INFO[0001] pulling sha256:cf8192ff388df08ea6446cc8e0046d5d33653108ada300a7436ce8dc8a0e4338
Hello world, I am wasi/wasm!
contents of / :
  3	.
  3	..
  4	hello

writing to /foo
contents of / after write :
  3	.
  3	..
  4	hello
  4	foo
```