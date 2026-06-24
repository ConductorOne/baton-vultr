package connector

import (
	"context"
	"io"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-vultr/pkg/client"
	cfg "github.com/conductorone/baton-vultr/pkg/config"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type Connector struct {
	client *client.VultrClient
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client),
		newACLbuilder(d.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (d *Connector) Asset(_ context.Context, _ *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(_ context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Baton Vultr connector",
		Description: "Implementation of the Vultr connector, with the resources users and ACLs.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(_ context.Context) (annotations.Annotations, error) {
	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, vultrBearerToken string) (*Connector, error) {
	l := ctxzap.Extract(ctx)

	vultrClient, err := client.New(ctx, vultrBearerToken)
	if err != nil {
		l.Error("error creating Vultr client", zap.Error(err))
		return nil, err
	}

	return &Connector{
		client: vultrClient,
	}, nil
}

// NewLambdaConnector returns a new ConnectorBuilderV2 for use with RunConnector.
func NewLambdaConnector(ctx context.Context, ac *cfg.Vultr, _ *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	c, err := New(ctx, ac.BearerToken)
	if err != nil {
		return nil, nil, err
	}
	return c, nil, nil
}
