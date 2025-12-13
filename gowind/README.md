# GoWind

## How To Use

### Install GoWind CLI

```shell
go install github.com/tx7do/kratos-cli/gowind@latest
```

### Create a New Project

```shell
gowind new myproject
cd myproject
go mod tidy
```

or you can use `-m` to specify the module name:

```shell
gowind new myproject -m github.com/yourusername/myproject
cd myproject
go mod tidy
```

### Run the microservice Application

```shell
gowind run
```

```shell
gowind run admin
```
