
# Docker Registry Manager

| Service   |  Master  | Develop  |   
|---|---|---|
| Status   | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master)  | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=develop)   |
| Coverage  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/develop/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  |
| Documentation  | [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  | [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  |

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many).

![Example](https://github.com/snagles/resources/blob/master/docker-registry-manager-updated.gif)

WARNING: This application is very much still a work in progress. Core functionality exists, but polish and features are still being worked on.

## Quickstart
 The below steps assume you have a docker registry currently running (with delete mode enabled (https://docs.docker.com/registry/configuration/). To add a registry to manage, add via the interface... or via the config.yml file

### Docker-Compose (Recommended)
 Install compose (https://docs.docker.com/compose/install/)

```bash
 git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
 vim config.yml # add your registry
 docker-compose up -d
 firefox localhost:8080
```

### Go
 ```bash
    git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
    vim config.yml # add your registry
    cd app && go build . && ./app
    firefox localhost:8080
 ```

### Dockerfile
 ```bash
    vim config.yml # add your registry
    docker run --detach --name docker-registry-manager -p 8080:8080 docker-registry-manager
 ```

## Current Features
 1. Support for docker distribution registry v2 (https and http)
 2. Viewable image/tags stages, commands, and sizes.
 3. Bulk deletes of tags
 4. Registry activity logs
 5. Comparison of registry images to Dockerhub public images

## Planned Features
 1. Authentication for users with admin/read only rights using TLS
 2. Global search
 3. List image shared layers
 4. Event timeline
