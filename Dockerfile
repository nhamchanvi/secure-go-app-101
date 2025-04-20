# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY *.go ./

# Build the static binary
# CGO_ENABLED=0 builds a static binary without C dependencies
# -ldflags="-w -s" strips debug information, making the binary smaller
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /main .

# Stage 2: Create the final minimal image
FROM alpine:latest 
# FROM scratch # Even smaller, but might need CA certs if making HTTPS calls

WORKDIR /app

# Create a non-root user and group
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

# Copy the static binary from the builder stage
COPY --from=builder /main /app/main

# Copy CA certificates if needed (especially if using scratch or making HTTPS calls)
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Ensure the binary is executable
RUN chmod +x /app/main

# Switch to the non-root user
USER nonroot

# Expose the port the app listens on
EXPOSE 8080

# Command to run the executable
CMD ["/app/main"]