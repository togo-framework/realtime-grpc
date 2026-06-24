// Package grpcrt is a gRPC transport for togo realtime. It implements togo.Broker:
// Publish fans events to in-process subscribers (the SSE Handler, for browsers) AND
// to connected gRPC streaming clients (for service-to-service realtime). A gRPC
// server runs on GRPC_ADDR exposing a server-streaming Subscribe method that uses a
// JSON codec, so no protoc-generated code is required.
//
//	togo install togo-framework/realtime-grpc
//
// Env: GRPC_ADDR (default :50051). gRPC clients dial with a JSON content-subtype and
// call togo.realtime.Realtime/Subscribe to receive {"event","data"} messages.
package grpcrt

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"

	"github.com/togo-framework/togo"
)

// Event is the wire message streamed to subscribers (JSON-encoded by the codec).
type Event struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

var hub = &broker{subs: map[chan Event]struct{}{}}

func init() {
	encoding.RegisterCodec(jsonCodec{})
	togo.RegisterProviderFunc("realtime-grpc", togo.PriorityService+10, func(k *togo.Kernel) error {
		addr := os.Getenv("GRPC_ADDR")
		if addr == "" {
			addr = ":50051"
		}
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			if k.Log != nil {
				k.Log.Warn("realtime-grpc: listen %s failed (%v); keeping default broker", addr, err)
			}
			return nil
		}
		srv := grpc.NewServer(grpc.ForceServerCodec(jsonCodec{}))
		srv.RegisterService(&serviceDesc, nil)
		go func() { _ = srv.Serve(lis) }()
		k.Realtime = hub
		return nil
	})
}

// ── Broker (togo.Broker) ────────────────────────────────────────────────────────

type broker struct {
	mu   sync.Mutex
	subs map[chan Event]struct{}
}

func (b *broker) subscribe() chan Event {
	ch := make(chan Event, 32)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *broker) unsubscribe(ch chan Event) {
	b.mu.Lock()
	delete(b.subs, ch)
	b.mu.Unlock()
}

func (b *broker) Publish(ev, data string) {
	e := Event{Event: ev, Data: data}
	b.mu.Lock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
	b.mu.Unlock()
}

func (b *broker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		ch := b.subscribe()
		defer b.unsubscribe(ch)
		fmt.Fprint(w, ": connected\n\n")
		flusher.Flush()
		for {
			select {
			case <-r.Context().Done():
				return
			case e := <-ch:
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Event, e.Data)
				flusher.Flush()
			}
		}
	}
}

// ── gRPC server (codegen-free via a JSON codec + manual ServiceDesc) ────────────

type jsonCodec struct{}

func (jsonCodec) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (jsonCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }
func (jsonCodec) Name() string                       { return "json" }

type realtimeServer any

var serviceDesc = grpc.ServiceDesc{
	ServiceName: "togo.realtime.Realtime",
	HandlerType: (*realtimeServer)(nil),
	Streams: []grpc.StreamDesc{
		{StreamName: "Subscribe", Handler: subscribeHandler, ServerStreams: true},
	},
	Metadata: "togo/realtime",
}

func subscribeHandler(_ any, stream grpc.ServerStream) error {
	var req map[string]any
	_ = stream.RecvMsg(&req) // one (ignored) request opens the stream
	ch := hub.subscribe()
	defer hub.unsubscribe(ch)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case e := <-ch:
			if err := stream.SendMsg(&e); err != nil {
				return err
			}
		}
	}
}
