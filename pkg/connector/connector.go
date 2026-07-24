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
	client        *client.VultrClient
	skipACLGrants bool
}

// ResourceSyncers returns a ResourceSyncerV2 for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(_ context.Context) []connectorbuilder.ResourceSyncerV2 {
	return []connectorbuilder.ResourceSyncerV2{
		newUserBuilder(d.client, d.skipACLGrants),
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

// New returns a new instance of the connector. syncACLs reports whether the acl
// resource type will be synced under the current configuration (derived from
// cli.ConnectorOpts.WillSyncResourceType in NewLambdaConnector); when false, the
// user syncer's resource type is annotated to skip both entitlements and grants
// for acl, since acl grants are emitted from the user syncer as an optimization.
func New(ctx context.Context, vultrBearerToken string, syncACLs bool) (*Connector, error) {
	l := ctxzap.Extract(ctx)

	vultrClient, err := client.New(ctx, vultrBearerToken)
	if err != nil {
		l.Error("error creating Vultr client", zap.Error(err))
		return nil, err
	}

	return &Connector{
		client:        vultrClient,
		skipACLGrants: !syncACLs,
	}, nil
}

// NewLambdaConnector returns a new ConnectorBuilderV2 for use with RunConnector.
func NewLambdaConnector(ctx context.Context, ac *cfg.Vultr, opts *cli.ConnectorOpts) (connectorbuilder.ConnectorBuilderV2, []connectorbuilder.Opt, error) {
	connectorOpts := opts
	if connectorOpts == nil {
		connectorOpts = &cli.ConnectorOpts{}
	}
	syncACLs := connectorOpts.WillSyncResourceType(AclResourceTypeID)

	c, err := New(ctx, ac.BearerToken, syncACLs)
	if err != nil {
		return nil, nil, err
	}
	return c, nil, nil
}
