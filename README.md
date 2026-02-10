![Baton Logo](./baton-logo.png)

# `baton-vultr` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-vultr.svg)](https://pkg.go.dev/github.com/conductorone/baton-vultr) ![ci](https://github.com/conductorone/baton-vultr/actions/workflows/ci.yaml/badge.svg)

`baton-vultr` is a connector for [Vultr](https://www.vultr.com/) built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

## Prerequisites

1. **Create a Vultr account**: If you don't have a **Vultr** account yet, you can register at [Vultr](https://www.vultr.com/).

2. **Get the Vultr API Key**:
    - Log in to your **Vultr** account.
    - Navigate to the **API** section from the control panel.
    - There you will be able to generate an **API Key** that will be used to authenticate Vultr API requests.

3. **Configure the API Key**:
    - You must provide your Vultr **API Key** through an environment variable or by passing it directly as a parameter when running tests.
    - You can set the `bearerToken` environment variable in your terminal with the following command:

      ````bash
      export bearerToken="your_api_key_here”
      ```
      **Note:** If you are using a development environment such as IntelliJ, you can also set this environment variable within the project run configuration.

## Connector capabilities

1. This connector synchronizes the resources users and ACLs.

2. This connector does not provision

## Connector credentials

1. This connector requires an API key that must be obtained when configuring the api on the vultr page [Vultr](https://www.vultr.com/).
   The main account from the api tab will be able to request the activation of the api and its API key.

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-vultr
baton-vultr
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_DOMAIN_URL=domain_url -e BATON_API_KEY=apiKey -e BATON_USERNAME=username ghcr.io/conductorone/baton-vultr:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-vultr/cmd/baton-vultr@main

baton-vultr

baton resources
```

# Data Model

`baton-vultr` will pull down information about the following resources:
- Users

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-vultr` Command Line Usage

```
baton-vultr

Usage:
  baton-vultr [flags]
  baton-vultr [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --bearer-token                 The client secret token used to authenticate with ConductorOne
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-vultr
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-vultr

Use "baton-vultr [command] --help" for more information about a command.
```
