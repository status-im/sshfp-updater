# syntax=docker/dockerfile:1
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod vendor
RUN go build -mod vendor  -o ./infra-sshfp-cf


FROM alpine:3.16.0
ARG REVISION=bf6021c8bb34394a70ed49c7e816b0aee4140992

LABEL org.opencontainers.image.authors="artur@status.im"
LABEL org.opencontainers.image.source="https://github.com/status-im/sshfp-generator"
LABEL org.opencontainers.image.revision=${REVISION}

WORKDIR /root
COPY --from=builder /app/infra-sshfp-cf ./

ENTRYPOINT [ "./infra-sshfp-cf" ]
