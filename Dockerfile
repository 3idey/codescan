# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE}" \
    -o /codescan ./cmd/codescan

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates git

# Create non-root user
RUN adduser -D -g '' codescan
USER codescan

WORKDIR /workspace

COPY --from=builder /codescan /usr/local/bin/codescan

ENTRYPOINT ["codescan"]
CMD ["--help"]
