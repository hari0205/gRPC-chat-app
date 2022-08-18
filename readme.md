
# gRPC Chat Application

A gRPC implementation of a chat application.


## Instruction

- Run the following to install protocol compiler plugins
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

```

- Update GO ENV path:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"

``` 

- Run the server

```bash
go run server.go
```
- Run client 
```bash
go run client/client.go

```

- In another terminal on the same directory , run another instance of the terminal