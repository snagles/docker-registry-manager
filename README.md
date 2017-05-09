
| Service   |  Master  | Develop  |   
|---|---|---|
| Status   | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master)  | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=develop)   |
| Coverage  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/develop/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  |
| Documentation  | [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  | [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  | 

# Docker Registry Manager

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many).

![Example](https://github.com/snagles/resources/blob/master/docker-registry-manager.gif)

WARNING: This application is very much still a work in progress. Core functionality exists, but polish and features are still being worked on.

## Quickstart
 The below steps assume you have a docker registry currently running (with delete mode enabled (https://docs.docker.com/registry/configuration/).

### Docker-Compose (Recommended)
 Install compose (https://docs.docker.com/compose/install/)

```bash
 git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
 docker-compose up -d
 firefox localhost:8080
```

### Go
 ```bash
    git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
    cd app && go build . && ./app -verbosity 6
    firefox localhost:8080
 ```

### Dockerfile
 ```bash
    docker build -t docker-registry-manager .
    docker run --detach --name docker-registry-manager -p 8080:8080 docker-registry-manager
 ```

## Current Features
 1. Support for docker distribution registry v2 (https and http)
 2. Manage multiple registries with one web instance
 3. Viewable image/tags stages, commands, and sizes. Refreshed every 45s
 4. Bulk deletes of tags
 5. Admin panel with logs, request tracking, and configurable log levels

## Planned Features
 1. Authentication for users with admin/read only rights
 2. Registry event logs in "dashboard"
