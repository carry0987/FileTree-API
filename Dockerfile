# Use an official Golang runtime as a parent image
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy the file from your host to your current location
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o filetree cmd/server/main.go

# Use a minimal alpine image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/filetree .

# Set the entry point
CMD ["./filetree"]
