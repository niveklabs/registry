# Build
FROM golang:1.14 AS builder

ARG VERSION

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -ldflags "-X main.version=${VERSION}" -o registry .

# # Final
FROM scratch

COPY --from=builder /build/registry /
COPY --from=builder /build/registry.json /

EXPOSE 8080

ENTRYPOINT [ "/registry" ]
