package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	echov1 "github.com/takekazu/golang-azd/pkg/gen/buf/echo/v1"        // generated by protoc-gen-go
	"github.com/takekazu/golang-azd/pkg/gen/buf/echo/v1/echov1connect" // generated by protoc-gen-connect-go
)

type EchoServer struct{}

func (s *EchoServer) Echo(
	ctx context.Context,
	req *connect.Request[echov1.EchoRequest],
) (*connect.Response[echov1.EchoResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&echov1.EchoResponse{
		Message: req.Msg.Message,
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	server := &EchoServer{}
	mux := http.NewServeMux()
	path, handler := echov1connect.NewEchoServiceHandler(
		server,
		connect.WithInterceptors(
		// https://github.com/takekazuomi/otel-azure-monitor/blob/main/cmd/server/server.go
		//			otelconnect.NewInterceptor(),
		))
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

// // newInterceptor instruments Connect clients and handlers using custom OpenTelemetry metrics, tracing, and propagation.
// func newInterceptor(tp trace.TracerProvider, mp metric.MeterProvider, p propagation.TextMapPropagator) (connect.Interceptor, error) {
// 	return otelconnect.NewInterceptor(
// 		otelconnect.WithTracerProvider(tp),
// 		otelconnect.WithMeterProvider(mp),
// 		otelconnect.WithPropagator(p),
// 	)
// }
