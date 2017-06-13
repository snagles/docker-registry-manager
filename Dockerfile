# Start from an Alpine image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# Install git
RUN apk add --no-cache git

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/snagles/docker-registry-manager/

# Install dependencies
WORKDIR $GOPATH/src/github.com/snagles/docker-registry-manager
RUN go get -v ./...

# Build the application inside the container.
WORKDIR $GOPATH/src/github.com/snagles/docker-registry-manager/app
RUN go build .

# Defaults
ENV MANAGER_PORT=8080
ENV MANAGER_LOG_LEVEL=5
ENV MANAGER_REFRESH_RATE=1m

# Run the app by default when the container starts.
CMD /go/src/github.com/snagles/docker-registry-manager/app/app
