FROM --platform=$TARGETPLATFORM golang:1.13.5-stretch as devel
COPY / /go/src/
RUN cd /go/src/cmd/opc && make all

FROM --platform=$TARGETPLATFORM busybox
COPY --from=devel /go/src/cmd/opc/baetyl-opc /bin/
ENTRYPOINT ["baetyl-opc"]
