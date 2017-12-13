# Start from an Alpine image with the latest version of Go installed
FROM golang:alpine

# Install git
RUN apk update && apk add git

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/snagles/docker-registry-manager/

# Build the application inside the container
WORKDIR $GOPATH/src/github.com/snagles/docker-registry-manager/app
RUN go build .

# Run the app by default when the container starts
CMD /go/src/github.com/snagles/docker-registry-manager/app/app
