############################
# STEP 1: Build executable binary
############################
# Use an official Go runtime as a base image
FROM golang:alpine AS builder

# Install git (required for fetching dependencies)
RUN apk update && apk add --no-cache git ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go service
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/kaze/

############################
# STEP 2: Build a small image
############################
FROM scratch

# Copy the CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app /app

# Copy the migrations directory (if exists and is required)
COPY --from=builder /app/migrations /migrations

# Expose the port that the application listens on
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app"]
