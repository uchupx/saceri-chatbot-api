
# --- Builder Stage ---
# Use a specific Go version that matches your go.mod file for consistency.
FROM golang:1.24-alpine AS builder

# Install necessary tools for building: make, git, and protobuf compiler.
RUN apk add --no-cache make git protobuf-dev

# Set the working directory inside the container.
WORKDIR /app

# Copy the Makefile first, as it might be needed for dependency steps.
COPY Makefile ./

# Copy go.mod and go.sum to leverage Docker's layer caching.
# This step will only be re-run if these files change.
COPY go.mod go.sum ./

# Download Go modules. Using `go mod download` is often more efficient in CI/CD
# than `go mod vendor`. If your `make vendor` does more, you can replace this.
RUN go mod download



# Install protobuf-related Go tools.
# These are required for `make proto` to generate Go code from .proto files.
RUN   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


# Copy the rest of the application source code.
COPY . .

# Generate protobuf Go files.
# This assumes your Makefile has a 'proto' target that runs protoc.
RUN make proto

# Build the Go application binary.
# The -o flag specifies the output file name.
# CGO_ENABLED=0 and -ldflags="-s -w" create a static, smaller binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o /saceri-chatbot-api ./cmd/api/main.go


# --- Final Stage ---
# Use a minimal base image for the final container.
FROM alpine:latest

# Install CA certificates for making HTTPS requests.
RUN apk --no-cache add ca-certificates

# Set the working directory.
WORKDIR /root/

# Copy the built binary from the builder stage.
COPY --from=builder /saceri-chatbot-api .

# Expose the port the application will run on.
# This should match the APP_PORT environment variable you provide at runtime.
EXPOSE 8080

# Command to run the application.
# The application will be started when the container launches.
CMD ["./saceri-chatbot-api"]
