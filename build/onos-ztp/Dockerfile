ARG ONOS_ZTP_BASE_VERSION=latest

FROM onosproject/onos-ztp-base:$ONOS_ZTP_BASE_VERSION as base

FROM alpine:3.9
RUN apk add libc6-compat

USER nobody

COPY --from=base /go/src/github.com/onosproject/onos-ztp/build/_output/onos-ztp /usr/local/bin/onos-ztp

ENTRYPOINT ["onos-ztp"]
