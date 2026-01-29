package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	bearerTokenField = field.StringField(
		"bearer-token",
		field.WithDisplayName("Bearer Token"),
		field.WithDescription("Bearer Token for authentication"),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)
	configurationFields = []field.SchemaField{
		bearerTokenField,
	}
	fieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run ./gen
var Config = field.Configuration{
	Fields:       configurationFields,
	Constraints:  fieldRelationships,
	DisplayName:  "Vultr",
	HelpUrl:      "/docs/baton/vultr",
	IconUrl:      "/static/app-icons/vultr.svg",
}
