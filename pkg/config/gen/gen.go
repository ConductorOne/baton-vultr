package main

import (
	"github.com/conductorone/baton-sdk/pkg/config"
	cfg "github.com/conductorone/baton-vultr/pkg/config"
)

func main() {
	config.Generate("vultr", cfg.Config)
}
