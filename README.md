# Go gRPC

My personal playing around with the proto family.

## Compile
I am committing the generated code. Still, the code to generate it is:

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto
```

## Env

```
docker run --name redis-test-instance -p 6379:6379 -d redis
```