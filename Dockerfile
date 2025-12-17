FROM golang:1.25.5-bookworm AS builder

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -v -o cert-manager-alidns-webhook

FROM debian:bookworm-slim

# update ca-certificates
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ARG user=cert-manager-alidns-webhook
ARG group=cert-manager-alidns-webhook
ARG uid=10000
ARG gid=10001

# If you bind mount a volume from the host or a data container,
# ensure you use the same uid
RUN groupadd -g ${gid} ${group} \
    && useradd -l -u ${uid} -g ${gid} -m -s /bin/bash ${user}

USER ${user}
WORKDIR /app

COPY --from=builder --chown=${uid}:${gid} /src/cert-manager-alidns-webhook /app/cert-manager-alidns-webhook

ENTRYPOINT ["/app/cert-manager-alidns-webhook"]
