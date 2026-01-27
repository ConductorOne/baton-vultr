package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	BearerTokenField = field.StringField(
		"bearer-token",
		field.WithDescription("Bearer Token for authentication"),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{BearerTokenField}

	// FieldRelationships defines relationships between the fields.
	FieldRelationships = []field.SchemaFieldRelationship{}

	// Config is the configuration schema for the connector.
	Config = field.Configuration{
		Fields:      ConfigurationFields,
		Constraints: FieldRelationships,
	}
)
