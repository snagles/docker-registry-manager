# Docker Registry Manager

Docker Registry Manager is a golang written, beego driven, web interface for interacting with multiple docker registries (one to many). Hopefully this application  can mature to the point ofan easy open source alternative to docker's trusted registry.

WARNING: This application is very much still a work in progress. Core functionality exists, but polish and features are still being worked on a daily basis. Heavy development is expected, so I encourage you to update regularly.

## Quickstart
 The below steps assume you have a docker registry currently running.
 
### Docker compose
 ```bash
    > git clone https://github.com/stefannaglee/docker-registry-manager.git
    > cd docker-registry-manager
    > vi docker-compose.yml # Edit the REGISTRYARGS line to include your registries in the format https://hostname:port/v2
    > docker-compose up
    > firefox localhost:8080 # for web ui
    > firefox localhost:8088 # for beego admin interface
 ```
 
### Go tool
 ```bash
    > git clone https://github.com/stefannaglee/docker-registry-manager.git
    > cd docker-registry-manager/app
    >  go build . && ./app -verbosity 6 -registry http://hostname:port/v2 # add more registries with another -registry flag
    > firefox localhost:8080 # for web ui
    > firefox localhost:8088 # for beego admin interface
 ```


## Current Features
 1. Support for docker distribution registry v2 (https and http).
 2. Support for multiple registries managed by one instance of this application
 3. Image/Tags stages, commands, and sizes
 3. Bulk deletes of repositories
 4. Viewable logs from the interface
 5. Admin interface using beego (on 8088) for tracking of request information
 
## Planned Features
 1. Docker compose support for multiple registries (that isn't a hack)
 2. Settings configuration
 3. Authentication for users with admin/read only rights etc.
 4. Authentication for registry access (using docker-registry auth)
 4. Activity log using the registries push events
 5. Dashboard for resource usage and other information
 6. Auotmated downloads of the latest public images for a repository stored on the registry of your choice
 7. Automated push and deployment to dockerhub on private registry push (if desired)
 8. Automated cleanup of images basic on configurable parameters
