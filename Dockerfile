# Use the official Golang image as the base image
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Change to the backend directory
WORKDIR /app/backend

# Build the Go application
RUN go build -o /app/HFTCryptoDashboard .

# Set the working directory for the final command
WORKDIR /app

# Command to run the application
CMD ["./HFTCryptoDashboard"]
