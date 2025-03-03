# API

This is the user and admin API for the Multi-Tenant Authentication service.

## Test & Build

To test and build the API, run the following command:

```bash
go test ./...
```

```bash
CGO_ENABLED=0 go build -ldflags=" \
    -X 'github.com/mmianl/api/v2/cmd/build.Version=0.0.0' \
    -X 'github.com/mmianl/api/v2/cmd/build.Commit=c30a9a606'" \
    -o out/api
```
