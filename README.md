# realtime-grpc

A **gRPC** transport for [togo](https://to-go.dev) realtime. Implements `togo.Broker`:
`Publish` fans events to in-process subscribers (the SSE `Handler`, for browsers) **and**
to connected gRPC streaming clients (service-to-service realtime). A gRPC server runs on
`GRPC_ADDR` exposing a server-streaming `Subscribe` method — using a JSON codec, so **no
protoc-generated code is required**.

## Install

```bash
togo install togo-framework/realtime-grpc
```

## Configure (`.env`)

```ini
GRPC_ADDR=:50051
```

gRPC clients dial with the `json` content-subtype and call
`togo.realtime.Realtime/Subscribe` to receive `{"event","data"}` messages.

MIT © togo-framework
