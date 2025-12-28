# Multi-stage build for Go application
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the spoticli-models dependency (for local replace directive)
COPY spoticli-models /spoticli-models

# Copy go mod files
COPY spoticli-backend/go.mod spoticli-backend/go.sum ./

RUN   ls -al

# Download dependencies
RUN go mod download

# Copy source code
COPY spoticli-backend/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spoticli-backend .

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/spoticli-backend .

# Change ownership
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 4200

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:4200/ || exit 1

# Run the application
CMD ["./spoticli-backend"]
