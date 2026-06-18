package config

import "github.com/conductorone/baton-sdk/pkg/field"

var (
	BearerTokenField = field.StringField(
		"bearer-token",
		field.WithDisplayName("Bearer Token"),
		field.WithPlaceholder("your-vultr-api-token"),
		field.WithDescription("Bearer Token for authentication"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	// ConfigurationFields defines the external configuration required for the connector to run.
	ConfigurationFields = []field.SchemaField{
		BearerTokenField,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConnectorDisplayName("Vultr"),
	field.WithIconUrl("/static/app-icons/vultr.svg"),
	field.WithHelpUrl("/docs/baton/vultr"),
)
