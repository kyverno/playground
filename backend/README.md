# Playground Backend

* Requires Go >= v1.20

## Install dependencies

```shell
go mod download
```

### Compiles backend

```shell
go build .
```

### Compiles and runs backend

```shell
go run .
```

You can pass flags when running the backend:

```shell
go run . --log=true --kubeconfig=~/.kube/config
```
