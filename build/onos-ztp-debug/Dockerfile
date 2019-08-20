ARG ONOS_ZTP_BASE_VERSION=latest

FROM onosproject/onos-ztp-base:$ONOS_ZTP_BASE_VERSION as base

FROM golang:1.12.6-alpine3.9 as debugBuilder

RUN apk upgrade --update --no-cache && apk add git && \
    go get -u github.com/go-delve/delve/cmd/dlv && \
    go get -u github.com/openconfig/gnmi/cmd/gnmi_cli

FROM alpine:3.9

RUN apk upgrade --update --no-cache && apk add bash bash-completion libc6-compat

COPY --from=base /go/src/github.com/onosproject/onos-ztp/build/_output/onos-ztp-debug /usr/local/bin/onos-ztp
COPY --from=debugBuilder /go/bin/dlv /usr/local/bin/dlv
COPY --from=debugBuilder /go/bin/gnmi_cli /usr/local/bin/gnmi_cli

RUN echo "#!/bin/sh" >> /usr/local/bin/onos-ztp-debug && \
    echo "dlv --listen=:40000 --headless=true --accept-multiclient=true --api-version=2 --log exec /usr/local/bin/onos-ztp -- \"\$@\"" >> /usr/local/bin/onos-ztp-debug && \
    chmod +x /usr/local/bin/onos-ztp-debug

RUN addgroup -S onos-ztp && adduser -S -G onos-ztp onos-ztp
USER onos-ztp
WORKDIR /home/onos-ztp

ENV ATOMIX_CONTROLLER=atomix-controller.kube-system.svc.cluster.local:5679
ENV ATOMIX_NAMESPACE=default

RUN cp /etc/profile /home/onos-ztp/.bashrc

ENTRYPOINT ["onos-ztp-debug"]
