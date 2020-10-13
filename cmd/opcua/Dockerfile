FROM --platform=$TARGETPLATFORM golang:1.13.5-stretch as devel
COPY / /go/src/
RUN cd /go/src/cmd/opcua && make all

FROM --platform=$TARGETPLATFORM busybox
COPY --from=devel /go/src/cmd/opcua/baetyl-opcua /bin/
ENTRYPOINT ["baetyl-opcua"]
