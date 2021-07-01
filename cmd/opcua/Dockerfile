FROM --platform=$TARGETPLATFORM golang:1.13.5-stretch as devel
COPY / /go/src/
RUN cd /go/src/cmd/opcua && make build-local BUILD_ARGS=$BUILD_ARGS

FROM --platform=$TARGETPLATFORM alpine
COPY --from=devel /go/src/cmd/opcua/opcua /bin/
ENTRYPOINT ["opcua"]
