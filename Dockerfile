# Go build stage - use bullseye for better CGO compatibility
# Use target platform for native compilation to support CGO dependencies
FROM golang:1.25-bookworm AS golang-builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Set working directory
WORKDIR /app

# Copy source code (excluding Dockerfile to avoid circular dependency)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build arguments for version info
ARG VERSION=docker
ARG COMMIT_HASH
ARG BUILD_TIME

# Build the streamable-http binary (CGO enabled for v8go dependency)
# Use native compilation instead of cross-compilation for CGO compatibility
RUN CGO_ENABLED=1 go build \
    -ldflags "-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildTime=${BUILD_TIME}" \
    -o resume-mcp-http \
    ./cmd/streamable-mcp/main.go


CMD ["./resume-mcp-http"]

# Final runtime stage
FROM ubuntu:24.04

COPY --from=golang-builder /app/resume-mcp-http /app/

# Install ca-certificates for HTTPS requests and wget for health check
RUN apt-get update && apt-get install -y ca-certificates tzdata wget && \
    rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -m appuser

# Set working directory
WORKDIR /app

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (default 8080, configurable via PORT env var)
EXPOSE 8080

# Environment variables with defaults
ENV PORT=8080
ENV GIN_MODE=release

# Run the binary
CMD ["./resume-mcp-http"]