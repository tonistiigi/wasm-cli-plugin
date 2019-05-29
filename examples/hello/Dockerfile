FROM --platform=$BUILDPLATFORM tonistiigi/xx:llvm AS build
ARG TARGETPLATFORM
WORKDIR /src
COPY hello.c .
RUN clang -static -o /hello hello.c

FROM scratch
ARG TARGETPLATFORM
ENV WHOAMI=$TARGETPLATFORM
COPY --from=build /hello .
ENTRYPOINT ["/hello"]