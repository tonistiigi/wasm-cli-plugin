# syntax=docker/dockerfile:1.1-experimental

ARG DOCKERD_VERSION=rc

FROM docker:$DOCKERD_VERSION AS dockerd-release

# xgo is a helper for golang cross-compilation
FROM --platform=$BUILDPLATFORM tonistiigi/xx:golang@sha256:6f7d999551dd471b58f70716754290495690efa8421e0a1fcf18eb11d0c0a537 AS xgo

FROM --platform=$BUILDPLATFORM golang:1.12-alpine AS gobase
COPY --from=xgo / /
RUN apk add --no-cache file git
ENV GOFLAGS=-mod=vendor
WORKDIR /src

FROM gobase AS wasm-version
RUN --mount=target=version,src=version \
  PKG=github.com/tonistiigi/wasm-cli-plugin VERSION=$(git describe --match 'v[0-9]*' --dirty='.m' --always --tags) REVISION=$(git rev-parse HEAD)$(if ! git diff --no-ext-diff --quiet --exit-code; then echo .m; fi); \
  echo "-X ${PKG}/version.Version=${VERSION} -X ${PKG}/version.Revision=${REVISION} -X ${PKG}/version.Package=${PKG}" | tee /tmp/.ldflags; \
  echo -n "${VERSION}" | tee /tmp/.version;

FROM gobase AS wasm-build
ENV CGO_ENABLED=0
ARG TARGETPLATFORM
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
  --mount=target=/go/pkg/mod,type=cache \
  --mount=source=/tmp/.ldflags,target=/tmp/.ldflags,from=wasm-version \
  set -x; go build -ldflags "$(cat /tmp/.ldflags)" -o /usr/bin/docker-wasm ./cmd/docker-wasm && \
  file /usr/bin/docker-wasm && file /usr/bin/docker-wasm | egrep "statically linked|Mach-O|Windows"

FROM wasm-build AS integration-tests
COPY . .

FROM --platform=$TARGETPLATFORM tonistiigi/wasmtime:binary AS wasmtime

FROM scratch AS binaries-unix
COPY --from=wasmtime / /
COPY --from=wasm-build /usr/bin/docker-wasm /

FROM binaries-unix AS binaries-darwin
FROM binaries-unix AS binaries-linux

FROM scratch AS binaries-windows
COPY --from=wasm-build /usr/bin/docker-wasm /docker-wasm.exe

FROM binaries-$TARGETOS AS binaries

FROM alpine AS demo-env
RUN apk add --no-cache iptables tmux ca-certificates curl
RUN mkdir -p /usr/local/lib/docker/cli-plugins && ln -s /usr/local/bin/docker-wasm /usr/local/lib/docker/cli-plugins/docker-wasm
RUN curl -fsSL https://tinyurl.com/install-buildx | sh
COPY ./hack/demo-env/entrypoint.sh /usr/local/bin
COPY ./hack/demo-env/tmux.conf /root/.tmux.conf
COPY --from=dockerd-release /usr/local/bin /usr/local/bin
WORKDIR /work
COPY ./examples .
COPY --from=binaries / /usr/local/bin/
VOLUME /var/lib/docker
ENV DOCKER_BUILDKIT=1
ENTRYPOINT ["entrypoint.sh"]

FROM alpine AS image
RUN apk add --no-cache ca-certificates
COPY --from=binaries / /bin/
ENTRYPOINT ["docker-wasm"]

FROM binaries