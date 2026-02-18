FROM golang:1.23-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on
RUN apk add --no-cache git
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /alloy-remote-config-server cmd/config/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /alloy-remote-config-server /alloy-remote-config-server
ENTRYPOINT ["/alloy-remote-config-server"]
