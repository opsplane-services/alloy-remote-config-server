package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	http "net/http"

	connect "connectrpc.com/connect"
	v1 "github.com/grafana/alloy-remote-config/api/gen/proto/go/collector/v1"
	collectorv1 "github.com/grafana/alloy-remote-config/api/gen/proto/go/collector/v1/collectorv1connect"
)

type Metadata struct {
	Id         string
	Attributes map[string]string
}

type ImplementedCollectorServiceHandler struct{}

func (ImplementedCollectorServiceHandler) GetConfig(
	ctx context.Context,
	req *connect.Request[v1.GetConfigRequest],
) (*connect.Response[v1.GetConfigResponse], error) {
	configID := req.Msg.GetId()
	attributes := req.Msg.GetLocalAttributes()
	metadata := Metadata{Id: configID, Attributes: attributes}
	templateName, ok := attributes["template"]
	if !ok {
		templateName = "default"
	}
	tmpl, ok := templates[templateName]
	if !ok {
		return nil, fmt.Errorf("Template %s not found", templateName)
	}
	var resolvedConfig strings.Builder
	err := tmpl.Execute(&resolvedConfig, metadata)
	if err != nil {
		return nil, err
	}
	globalStorage.Set(configID, resolvedConfig.String())
	res := connect.NewResponse(&v1.GetConfigResponse{Content: resolvedConfig.String()})
	return res, nil
}

func (ImplementedCollectorServiceHandler) RegisterCollector(
	ctx context.Context,
	req *connect.Request[v1.RegisterCollectorRequest],
) (*connect.Response[v1.RegisterCollectorResponse], error) {
	configID := req.Msg.GetId()
	log.Printf("Register: %v [not used - agents are registered by getConfig call]", configID)
	res := connect.NewResponse(&v1.RegisterCollectorResponse{})
	return res, nil
}

func (ImplementedCollectorServiceHandler) UnregisterCollector(
	ctx context.Context,
	req *connect.Request[v1.UnregisterCollectorRequest],
) (*connect.Response[v1.UnregisterCollectorResponse], error) {
	configID := req.Msg.GetId()
	log.Printf("Unregister: %v [not used - agents are unregistered once not accessed for long time]", configID)
	res := connect.NewResponse(&v1.UnregisterCollectorResponse{})
	return res, nil
}

func StartConnectGrpcServer(listenAddr string, port int) {
	mux := http.NewServeMux()
	mux.Handle(collectorv1.NewCollectorServiceHandler(&ImplementedCollectorServiceHandler{}))
	log.Printf("Start listening (gRPC) on port %d", port)
	err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", listenAddr, port),
		h2c.NewHandler(mux, &http2.Server{}),
	)
	log.Fatalf("listen failed: %v", err)
}
