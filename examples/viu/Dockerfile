FROM --platform=$BUILDPLATFORM tonistiigi/xx:rust AS build
RUN apt-get update && apt-get install -y git
WORKDIR /src
RUN git clone git://github.com/wapm-packages/viu
WORKDIR /src/viu
ARG TARGETPLATFORM
RUN cargo build --release --out-dir=/out

FROM debian:buster-slim AS release-linux
COPY --from=build /out/ /usr/bin/
ENTRYPOINT ["viu"]

FROM scratch AS release-wasi
COPY --from=build /out/viu.wasm /
ENTRYPOINT ["/viu.wasm"]

FROM release-$TARGETOS AS release
ADD https://media.giphy.com/media/xT0xeJpnrWC4XWblEk/giphy.gif /mindblown.gif
ADD docker.png /

FROM release AS docker
CMD ["/docker.png"]

FROM release
