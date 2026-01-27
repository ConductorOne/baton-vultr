package config

//go:generate go run ./gen

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var Config = field.NewConfiguration([]field.SchemaField{
	field.StringField(
		"bearer-token",
		field.WithDescription("Bearer Token for authentication"),
		field.WithRequired(true),
	),
})

func ValidateConfig(c *Vultr) error {
	return nil
}
