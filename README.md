# Utility to use for e2guardian_setup

## How to use
```shell
go build utils.go
```

### Generate X25519Identity
./utils -generate_X25519_key="true"

### Encode some text to base64
./utils -to_encode="some text to encode as base64"
