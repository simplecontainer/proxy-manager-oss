FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:latest as initial

ARG BINARY=proxy-manager-oss
ARG DOCKER_BINARY=proxy-manager-oss
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN adduser -D node && \
    mkdir -p /opt/app /home/node/app /home/node/.ssh/ && \
    chown -R node:node /home/node /opt/app && \
    touch /home/node/.ssh/known_hosts

COPY $BINARY-$TARGETOS-$TARGETARCH/$DOCKER_BINARY /opt/app/$DOCKER_BINARY
RUN chmod +x /opt/app/$DOCKER_BINARY

RUN apk add --no-cache busybox-extras

FROM scratch

COPY --from=initial / /

USER node
ENTRYPOINT ["/opt/app/proxy-manager-oss"]