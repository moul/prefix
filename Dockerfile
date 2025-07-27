# Build stage
FROM golang:1.23-alpine as builder

RUN apk add --no-cache git gcc musl-dev make

WORKDIR /build
COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o prefix ./cmd/prefix

# Runtime stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/prefix /usr/local/bin/prefix

ENTRYPOINT ["prefix"]