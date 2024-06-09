package internal

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"
)

var tmplFuncs = map[string]any{
	"rpcMethod":   rpcMethod,
	"methodInput": methodInput,
}

// FileTmpl is .jgw template
var FileTmpl = template.Must(template.New("").Funcs(tmplFuncs).Parse(`
// Code generated by protoc-gen-jrpc-gateway. DO NOT EDIT.
// source: {{ .Name }}

{{$package := .Package}}
/*

Package {{ $package }} is a reverse proxy.

It translates gRPC into JSON-RPC 2.0
*/
package {{ $package }}

import (
	"context"
	"encoding/json"

    "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

{{range $service := .Service}}

{{$serviceName := $service.GetName | printf "%sJsonRpcService"}}
{{$clientName := $service.GetName | printf "%sClient"}}
{{$clientConstructor := printf "New%sClient" $service.GetName}}

type {{$serviceName}} struct {
	client	{{$clientName}}
}

type paramsAndHeaders struct {
	Headers metadata.MD ` + "`json:\"headers,omitempty\"`" + `
	Params  json.RawMessage   ` + "`json:\"params\"`" + `
}

// {{$serviceName | printf "Register%s"}} register the grpc client {{$service.GetName}} for json-rpc.
// The handlers forward requests to the grpc endpoint over "conn".
func {{$serviceName | printf "Register%s"}} (conn *grpc.ClientConn) *{{$serviceName}} {
	return &{{$serviceName}} {
		client: {{$clientConstructor}}(conn),
	}
}

func (s *{{$serviceName}}) Methods() map[string]func(ctx context.Context, message json.RawMessage) (any, error) {
	return map[string]func(ctx context.Context, params json.RawMessage) (any, error) {
		{{range $method := $service.GetMethod}}
			"{{rpcMethod $package $service.GetName $method.GetName}}": func(ctx context.Context, data json.RawMessage) (any, error) {
				req := new({{methodInput $method.GetInputType}})

				var jrpcData paramsAndHeaders
				
				if err := json.Unmarshal(data, &jrpcData); err != nil {
					return nil, err
               }

				err := protojson.Unmarshal(jrpcData.Params, req)
				if err != nil {
					return nil, err
				}

				return s.client.{{$method.GetName}}(metadata.NewOutgoingContext(ctx, jrpcData.Headers), req)
			},
		{{end}}
	}
}

{{end}}
`))

func rpcMethod(pkg, service, method string) string {
	return fmt.Sprintf("%s.%s.%s", camelToSnake(pkg), camelToSnake(service), camelToSnake(method))
}

func camelToSnake(s string) string {
	var buf bytes.Buffer
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteRune(unicode.ToLower(r))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func methodInput(in string) string {
	sep := strings.Split(in, ".")
	return sep[len(sep)-1]
}
