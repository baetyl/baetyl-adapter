FROM --platform=$TARGETPLATFORM golang:1.18.3-stretch as devel
ARG BUILD_ARGS
COPY / /go/src/
RUN cd /go/src/cmd/modbus && make build-local BUILD_ARGS=$BUILD_ARGS

FROM --platform=$TARGETPLATFORM alpine
COPY --from=devel /go/src/cmd/modbus/modbus /bin/
ENTRYPOINT ["modbus"]
