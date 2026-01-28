package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	cfg "github.com/conductorone/baton-vultr/pkg/config"
	"github.com/conductorone/baton-vultr/pkg/connector"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-vultr",
		getConnector,
		cfg.Config,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.Connector{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, vc *cfg.Vultr) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	if err := cfg.ValidateConfig(vc); err != nil {
		return nil, err
	}

	connectorBuilder, err := connector.New(ctx, vc.BearerToken)
	if err != nil {
		l.Error("error creating newConnector", zap.Error(err))
		return nil, err
	}

	opts := make([]connectorbuilder.Opt, 0)

	newConnector, err := connectorbuilder.NewConnector(ctx, connectorBuilder, opts...)
	if err != nil {
		l.Error("error creating newConnector", zap.Error(err))
		return nil, err
	}
	return newConnector, nil
}
