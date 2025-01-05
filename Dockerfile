
# Use the official Golang image as a base for building the Go application
FROM golang:1.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Set the GOARCH environment variable for cross-compiling (for amd64)
ENV GOARCH=amd64

# Build the Go application for amd64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jadwal_sholat .

# Use a minimal image for running the compiled binary (e.g., Debian)
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/jadwal_sholat .

# Copy the JSON file
COPY jadwal_sholat.json .

# Copy the url text
COPY url.text .

# Set the command to run the application
CMD ["./jadwal_sholat"]
