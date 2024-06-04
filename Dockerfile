# Use the official Golang image as the base image
FROM golang:1.21.6

# Set the working directory inside the container
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

# Copy the local main files to the container's workspace
COPY ./ ${ROOT}

# Copy the Go module files
COPY go.mod go.sum ${ROOT}

# Download and install the Go dependencies
RUN go mod download
