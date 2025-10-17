# alloy-remote-config-server

![build workflow](https://github.com/opsplane-services/alloy-remote-config-server/actions/workflows/docker-publish.yml/badge.svg)
![license](http://img.shields.io/badge/license-Apache%20v2-blue.svg)

## Description

Remote server implementation of alloy remote config API.

See: [Grafana alloy remote config repo](https://github.com/grafana/alloy-remote-config)

### Features
- Centralized configuration management (by static template files - ideally use GitOps)
- RESTful API for fetching configurations and templates
- Keep resolved configuration (by host id) in Redis
- Easy to deploy with Docker

This implementation provides gRPC endpoint (and additional HTTP endpoint), that can be used by `remotecfg` block of Grafana alloy block. The server uses the `id` and `local_attributes` fields to fill a predefined go template file. (where those variables can be used for templating) A template folder can be pre-defined, all the templates in that folder that ends with `.conf.tmpl` suffix will be loaded into the application. 

At least a default template configuration (`default.conf.tmpl`) is required, and the tempaltes can be referenced in the `local_attributes` field by the names without the file extension suffix. (e.g.: `default` as `template` attribute from Grafana Alloy agent configuration - the proper template can be selected by this `template` attribute). 

The resolved configurations are also stored in the application in memory or in redis. (these can be accessed by `/configs` or `/configs/{:id}` http endpoints)

If you would like to use TLS / MTLS or OAuth 2.0 for this gRPC server implementation, It's recommend to deploy the server to Kubernetes and use something like Istio and set the proper request authentication or authorization policies around the service.

Note: this implementation uses `attributes` instead of `manifest` or `local_attributes` fields for the remotecfg. (to be up-to-date with the alloy implementation, but this is not up-to-date with its current documentation, probably `local_attributes` will be used soon)

## Configuration

The following environment variables can be used by the application (or set throgh `.env` file):

- `CONFIG_FOLDER`: Directory that should contain the static go template configuration files with `.conf.tmpl` extension. (default value: `conf`)
- `GRPC_PORT`: GRPC port for the config service that implements `GetConfig` operation. (default value: `8888`)
- `HTTP_PORT`: HTTP port for the additional web service to query templates and resolved configurations (default value: `8080`)
- `USE_REDIS`: (default value: `false`)
- `REDIS_URL`: Redis URL that is parsed at application startup if Redis is used (can contain username/password )
- `REDIS_TTL`: TTL value that is set for the resolved configuration objects - once a config with the same id is resolved again, the TTL is re-set (default value: `259200`)
- `ORG_NAME`: Organization name is a global configuration that can be used to separate configs based on this namespace - therefore you can run multiple instances of this application with different organization names - but using the same Redis storage (default value: `default`)
- `LISTEN_ADDR`: Host to listen by the GRPC/HTTP server. (default value: `0.0.0.0`)

## Usage

### Docker

```bash
docker pull opsplane/alloy-remote-config-server:latest
# use -e for setting environment variables or pass .env file though a volume with -v
docker run opsplane/alloy-remote-config-server:latest
```

### Local

```bash
# create .env file from .env.template
go mod tidy
go build cmd/config/main.go
go run cmd/config/main.go
```

## License

MIT license
