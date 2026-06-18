package main

import (
	cfg "github.com/conductorone/baton-vultr/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("vultr", cfg.Config)
}
