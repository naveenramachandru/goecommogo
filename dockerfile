# Use the official Golang image as a base image
FROM golang:1.2.0

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy your application source code to the container
COPY main.go ./

# Build your Go application
RUN go build -o bin .

# Set the entry point for your application
CMD ["./bin"]
