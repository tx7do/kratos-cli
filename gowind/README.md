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

### Run the microservice Application

You can directly execute the microservice in the current path without parameters, For example you are currently in 'app/admin/service':

```shell
gow run
```

or run a specified microservice, for example `admin` service:

```shell
gow run admin
```
