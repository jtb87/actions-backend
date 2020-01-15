# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.13-alpine base image
FROM golang:1.13-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade 
# && \
#     apk add --no-cache bash git 

# Set the Current Working Directory inside the container

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .
# Setup application

COPY . /opt/backend
WORKDIR /opt/backend

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Build the Go app
RUN go build -o main .

# Run the executable
CMD ["./main"]