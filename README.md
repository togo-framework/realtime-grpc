<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/realtime-grpc</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/realtime-grpc"><img src="https://pkg.go.dev/badge/github.com/togo-framework/realtime-grpc.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/realtime-grpc
```

<!-- /togo-header -->

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

<!-- togo-sponsors -->
---

<div align="center">
  <h3>Premium sponsors</h3>
  <p>
    <a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp;
    <a href="https://one-studio.co"><strong>One Studio</strong></a>
  </p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
