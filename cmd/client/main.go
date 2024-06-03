package main

import (
	"context"
	"log"
	"net/http"
	"os"

	echov1 "github.com/takekazu/golang-azd/pkg/gen/buf/echo/v1"
	"github.com/takekazu/golang-azd/pkg/gen/buf/echo/v1/echov1connect"
	"google.golang.org/protobuf/encoding/protojson"

	"connectrpc.com/connect"
)

func main() {
	baseURL := "http://localhost:8080"

	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	client := echov1connect.NewEchoServiceClient(
		http.DefaultClient,
		baseURL,
		connect.WithGRPC(),
	)
	res, err := client.Echo(
		context.Background(),
		connect.NewRequest(&echov1.EchoRequest{Message: "Hello, connectRPC"}),
	)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(protojson.Format(res.Msg))
}
