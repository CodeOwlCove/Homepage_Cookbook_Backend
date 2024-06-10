# Start from the official golang image
FROM golang:1.17-alpine

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create a directory inside the container to store all our application and then make it the working directory.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -tags prod -o main ./src/main

# Expose port 8085 to the outside world
EXPOSE 8085

# Command to run the executable
CMD ["./main"]