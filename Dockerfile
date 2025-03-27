# Start from the official Go image
FROM golang:1.24.0-alpine3.20 AS builder

# Set working directory inside the container
WORKDIR /go/src/github.com/schlafer/EventApp

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
# Adjust the path to your main.go file
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Use a minimal Alpine image for the runtime environment
FROM alpine:3.20

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory inside the container
WORKDIR /usr/bin

# Copy the pre-built binary from the previous stage
COPY --from=builder /go/src/github.com/schlafer/EventApp/main ./app

# Expose port (change to your API's port)
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
