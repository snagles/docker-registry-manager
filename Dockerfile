# Start from an Alpine image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# Install git
RUN apk add --no-cache git

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/stefannaglee/docker-registry-manager/

# Install dependencies
WORKDIR $GOPATH/src/github.com/stefannaglee/docker-registry-manager/

RUN go get -v -d ./...
RUN go get -v github.com/smartystreets/goconvey

# Build the application inside the container.
RUN cd app/ && go build

# Set env for verbose output (Default info level)
ENV VERBOSITY 5

# Run the app by default when the container starts.
CMD $GOPATH/src/github.com/stefannaglee/docker-registry-manager/app/app -verbosity $VERBOSITY $REGISTRYARGS
