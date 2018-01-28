
# Docker Registry Manager [![Go Report Card](https://goreportcard.com/badge/github.com/snagles/docker-registry-manager)](https://goreportcard.com/report/github.com/snagles/docker-registry-manager) [![GoDoc](https://godoc.org/github.com/snagles/docker-registry-manager?status.svg)](https://godoc.org/github.com/snagles/docker-registry-manager)  

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many).

| Service   |  Master  | Develop  |   
|---|---|---|
| Status   | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=master)  | ![Build Status](https://travis-ci.org/snagles/docker-registry-manager.svg?branch=develop)   |
| Coverage  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  | [![Coverage Status](https://codecov.io/gh/snagles/docker-registry-manager/branch/develop/graph/badge.svg)](https://codecov.io/gh/snagles/docker-registry-manager)  |

![Example](https://github.com/snagles/resources/blob/master/docker-registry-manager-updated.gif)

## Current Features
 1. Support for docker distribution registry v2 (https and http)
 2. Viewable image/tags stages, commands, and sizes.
 3. Bulk deletes of tags
 4. Registry activity logs
 5. Comparison of registry images to public Dockerhub images

## Planned Features
 1. Authentication for users with admin/read only rights using TLS
 2. Global search
 3. List image shared layers
 4. Event timeline

## Quickstart
 The below steps assume you have a docker registry currently running (with delete mode enabled (https://docs.docker.com/registry/configuration/). To add a registry to manage, add via the interface... or via the registries.yml file

### Docker-Compose (Recommended)
 Install compose (https://docs.docker.com/compose/install/), and then run the below commands

 ```bash
  git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
  vim registries.yml # add your registry
  vim docker-compose.yml # Edit application settings e.g log level, port
  docker-compose up -d
  firefox localhost:8080
  ```

#### Environment Options:
 - MANAGER_PORT: Port to run on inside the docker container
 - MANAGER_REGISTRIES: Registries.yml file location inside the docker container
 - MANAGER_LOG_LEVEL: Log level for logs (fatal, panic, error, warn, info, debug)
 - MANAGER_ENABLE_HTTPS: true/false for using HTTPS. When using HTTPS the below options must be set
 - MANAGER_KEY: key file location inside the docker container
 - MANAGER_CERTIFICATE: Certificate location inside the docker container

### Go
 ```bash
    git clone https://github.com/snagles/docker-registry-manager.git && cd docker-registry-manager
    vim registries.yml # add your registry
    cd app && go build . && ./app --port 8080 --log-level warn --registries "../registries.yml"
    firefox localhost:8080
 ```

#### CLI Options
  - port, p: Port to run on
  - registries, r: Registrys.yml file location
  - log-level, l: Log level for logs (fatal, panic, error, warn, info, debug)
  - enable-https, e: true/false for using HTTPS. When using HTTPS the below options must be set
  - tls-key, k: key file location inside the docker container
  - tls-certificate, cert: Certificate location inside the docker container

### Dockerfile
 ```bash
    vim registries.yml # add your registry
    docker run --detach --name docker-registry-manager -p 8080:8080 -e MANAGER_PORT=8080 -e MANAGER_REGISTRIES=/app/registries.yml -e MANAGER_LOG_LEVEL=warn docker-registry-manager
    firefox localhost:8080
 ```

#### Environment Options:
- MANAGER_PORT: Port to run on inside the docker container
- MANAGER_REGISTRIES: Registries.yml file location inside the docker container
- MANAGER_LOG_LEVEL: Log level for logs (fatal, panic, error, warn, info, debug)
- MANAGER_ENABLE_HTTPS: true/false for using HTTPS. When using HTTPS the below options must be set
- MANAGER_KEY: key file location inside the docker container
- MANAGER_CERTIFICATE: Certificate location inside the docker container

### Registries.yml Example
```yml
registries:
  localRegistry:
    url: http://localhost # Example https://localhost, http://remotehost.com
    port: 5000  # Example: 443, 8080, 5000
    username: exampleUser
    password: examplePassword
    refresh-rate: "5m" # Example: 60s, 5m, 1h
    skip-tls-validation: true # REQUIRED for self signed certificates
    dockerhub-integration: true # Optional - compares to dockerhub to determine if image up to date
```
