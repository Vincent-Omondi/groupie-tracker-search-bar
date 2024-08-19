# Use golang:1.22.2-alpine as the base image
FROM golang:1.22.2-alpine

# Additional image metadata
LABEL maintainer="github.com/johneliud"
LABEL version="1.0"
LABEL description="Image file for a web server that consists of receiving a given API and manipulate the data contained in it in order to create a site displaying the information."

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod file
COPY go.mod .

# Install go.mod dependancies
RUN go.mod download

# Copy the rest of the source code
COPY . .

# Build an executable file from the program
RUN go build -o main .

# Expose container to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]