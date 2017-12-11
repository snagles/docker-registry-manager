# Start from an Alpine image with the latest version of Go installed
FROM golang:alpine as build-env

# Install git
RUN apk add --no-cache git

WORKDIR $GOPATH/src/github.com/snagles/docker-registry-manager

# Copy the local package files to the container's workspace.
ADD . ./

# Build the application inside the container
RUN go install github.com/snagles/docker-registry-manager/app


# Distribution image
FROM alpine:3.7

# Copy binary from build stage
COPY --from=build-env /go/bin/app /app/docker-registry-manager

# Set the default config location
ENV REGISTRY_CONFIG /var/lib/docker-registry-manager
# Config storage volume for persistent settings
VOLUME ["/var/lib/docker-registry-manager"]

# Run the app by default when the container starts
CMD /app/docker-registry-manager
