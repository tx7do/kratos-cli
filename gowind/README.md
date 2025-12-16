# GoWind


## How To Use

### Install GoWind CLI

```shell
go install github.com/tx7do/kratos-cli/gowind/cmd/gow@latest
```

### Create a New Project

```shell
gow new myproject
cd myproject
go mod tidy
```

or you can use `-m` to specify the module name:

```shell
gow new myproject -m github.com/yourusername/myproject
cd myproject
go mod tidy
```

### Add a New Microservice

```shell
gow add service admin
gow add service user
go mod tidy
```

You can also specify the server type and orm type when adding a new microservice:

```shell
# add microservice with grpc
gow add service order -s grpc

# add microservice with rest
gow add service admin -s rest

# add microservice with rest and grpc servers
gow add service admin -s rest -s grpc

# add microservice with gorm orm and grpc server
gow add svc payment -d gorm -s grpc

# add microservice with rest, grpc servers and gorm, redis daos
gow add service admin -s rest -s grpc -d gorm -d redis
```

### Run The Microservice Application

You can directly execute the microservice in the current path without parameters, For example you are currently in '
app/admin/service':

```shell
gow run
```

or run a specified microservice, for example `admin` service:

```shell
gow run admin
```

### Ent Code Generation

add ent schema for a microservice, for example `admin` service:

```shell
gow ent add admin User,Group
```

generate ent code for all microservices:

```shell
gow ent generate
```

generate ent code for a specified microservice, for example `admin` service:

```shell
gow ent generate admin
```

## Wire Code Generation

generate wire code for all microservices:

```shell
gow wire generate
```

generate wire code for a specified microservice, for example `admin` service:

```shell
gow wire generate admin
```
