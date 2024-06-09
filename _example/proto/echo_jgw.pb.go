// Code generated by protoc-gen-jrpc-gateway. DO NOT EDIT.
// source: echo.proto

/*
Package proto is a reverse proxy.

It translates gRPC into JSON-RPC 2.0
*/
package proto

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

type EchoServiceJsonRPC struct {
	client EchoServiceClient
}

type paramsAndHeadersEchoService struct {
	Headers metadata.MD     `json:"headers,omitempty"`
	Params  json.RawMessage `json:"params"`
}

// RegisterEchoServiceJsonRPC register the grpc client EchoService for json-rpc.
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterEchoServiceJsonRPC(conn *grpc.ClientConn) *EchoServiceJsonRPC {
	return &EchoServiceJsonRPC{
		client: NewEchoServiceClient(conn),
	}
}

func (s *EchoServiceJsonRPC) Methods() map[string]func(ctx context.Context, message json.RawMessage) (any, error) {
	return map[string]func(ctx context.Context, params json.RawMessage) (any, error){

		"proto.echo_service.echo": func(ctx context.Context, data json.RawMessage) (any, error) {
			req := new(EchoRequest)

			var jrpcData paramsAndHeadersEchoService

			if err := json.Unmarshal(data, &jrpcData); err != nil {
				return nil, err
			}

			err := protojson.Unmarshal(jrpcData.Params, req)
			if err != nil {
				return nil, err
			}

			return s.client.Echo(metadata.NewOutgoingContext(ctx, jrpcData.Headers), req)
		},
	}
}
