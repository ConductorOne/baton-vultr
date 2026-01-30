package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	bearerTokenField = field.StringField(
		"bearer-token",
		field.WithDisplayName("Bearer Token"),
		field.WithDescription("Bearer Token for authentication"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the
	// connector to run.
	ConfigurationFields = []field.SchemaField{
		bearerTokenField,
	}

	// FieldRelationships defines relationships between the fields listed in
	// ConfigurationFields that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Vultr"),
	field.WithHelpUrl("/docs/baton/vultr"),
	field.WithIconUrl("/static/app-icons/vultr.svg"),
)
