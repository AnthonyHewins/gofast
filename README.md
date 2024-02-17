# GoFAST

CLI application to quickly deploy servers/CLIs with my opinion about what's considered good practice. Servers are built to be deployed in a k8s environment.

## Feature list

**App creation**

- Creates a CLI using `cobra-cli`
- Creates a gRPC/REST compatible server using
    - `grpc-gateway`: write gRPC, get an entire REST server for free
    - `buf`: package management for protobuf
    - Prometheus metrics
    - Kube health gRPC integration
- SQLc for typesafe SQL in the `sql` directory

## Install

```shell
go install github.com/AnthonyHewins/gofast@latest
```

## Create a new application

```shell
export USER=
export PROJECT=
gofast create github.com/${USER}/${PROJECT}
```