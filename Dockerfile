# Step 1: build image
FROM golang:1.21 AS builder

# Cache the dependencies
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

# Compile the application
COPY . /app
RUN ./scripts/build.sh

# Step 2: build the image to be actually run
FROM alpine:3.18.4
COPY --from=builder /app/bin/volume_injector /app/bin/volume_injector
ENTRYPOINT ["/app/bin/volume_injector"]