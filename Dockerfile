# Start from an Alpine image with the latest version of Go installed
FROM golang:alpine as build-env

# Install git and the bee tool used for deployment
RUN apk add --no-cache git

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/snagles/docker-registry-manager

# Build the application using the bee tool
RUN go get github.com/beego/bee
RUN bee pack -p /go/src/github.com/snagles/docker-registry-manager/app
RUN mkdir /app
RUN tar -xzvf /go/app.tar.gz --directory /app

# Distributed image
FROM alpine:3.7
RUN apk add --no-cache ca-certificates

# Copy packed beego tar
WORKDIR /app
COPY --from=build-env /app /app
COPY --from=build-env /go/src/github.com/snagles/docker-registry-manager/registries.yml /app/registries.yml

# Set the default registries location and volume
ENV MANAGER_REGISTRIES=/app/registries.yml
ENV MANAGER_LOG_LEVEL=warn
VOLUME ["/app"]

# Run the app by default when the container starts
CMD ["/app/app"]
