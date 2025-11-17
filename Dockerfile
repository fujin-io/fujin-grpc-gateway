ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && apk add git make

COPY . .

RUN make build

FROM scratch

WORKDIR /

COPY --from=builder app/bin/fujin-grpc-gateway /fujin-grpc-gateway

STOPSIGNAL SIGTERM

ENTRYPOINT ["/fujin-grpc-gateway"]