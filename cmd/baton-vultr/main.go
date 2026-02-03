package main

import (
	"context"

	cfg "github.com/conductorone/baton-vultr/pkg/config"
	"github.com/conductorone/baton-vultr/pkg/connector"
	"github.com/conductorone/baton-sdk/pkg/config"
)

var version = "dev"

func main() {
	ctx := context.Background()
	config.RunConnector(ctx,
		"baton-vultr",
		version,
		cfg.Config,
		connector.New,
	)
}
