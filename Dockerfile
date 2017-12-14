# Start from an Alpine image with the latest version of Go installed
FROM golang:alpine as build-env

# Install git and the bee tool used for deployment
RUN apk add --no-cache git

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/snagles/docker-registry-manager

# Build the application using the bee tool
RUN go get github.com/beego/bee
RUN bee pack -p /go/src/github.com/snagles/docker-registry-manager/app

# Distributed image
FROM alpine:3.7
RUN apk add --no-cache tar ca-certificates

# Copy packed beego tar
WORKDIR /app
COPY --from=build-env /go/app.tar.gz /app/app.tar.gz
RUN tar -xzf app.tar.gz \
    && rm app.tar.gz

# Set the default config location and volume
ENV REGISTRY_CONFIG /var/lib/docker-registry-manager/config.yml
VOLUME ["/var/lib/docker-registry-manager"]

# Run the app by default when the container starts
CMD /app/app
