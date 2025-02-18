# Stage 1: Build Go Application
FROM golang:1.23-alpine as base

# Install dependencies for building Go apps
RUN apk add --no-cache build-base git inotify-tools bash

# Set working directory
WORKDIR /clean_web

# Copy go.mod and go.sum separately to leverage Docker cache
COPY go.mod go.sum ./

# Initialize Go modules
RUN go mod download

# Copy all files from your repo to the container
COPY . .

# Build Go application
RUN go build -o main .

# Clean up build dependencies to reduce image size
RUN apk del build-base git && rm -rf /var/cache/apk/*


# Stage 2: Final Image with FFmpeg, Go, and espeak
FROM jrottenberg/ffmpeg:4.4-alpine

# Install dependencies
RUN apk add --no-cache bash inotify-tools git espeak python3 py3-pip

# Install gTTS
RUN pip install gtts

# Install Go manually (karena base image tidak punya Go)
RUN apk add --no-cache --virtual .build-deps curl && \
    curl -fsSL https://golang.org/dl/go1.23.0.linux-amd64.tar.gz | tar -C /usr/local -xz && \
    apk del .build-deps

# Set Go environment variables
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

# Copy the built Go binary from the first stage
COPY --from=base /clean_web/main /usr/local/bin/

# Set working directory for the final app
WORKDIR /clean_web

# Copy necessary files (run.sh, configs, etc.)
COPY . .

# Set execution permission
RUN chmod +x /usr/local/bin/main

# Default command to run the application
ENTRYPOINT ["/usr/local/bin/main"]
