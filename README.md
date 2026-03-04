# ha-ctl

A CLI tool for controlling Home Assistant, designed for use as an LLM agent skill.

## Features

- **Entity search** — find entities by name with `ha-ctl find`
- **Live state** — get real-time entity state with `ha-ctl state`
- **Service calls** — control devices with `ha-ctl call`
- **Home context** — compact home summary optimized for LLM token budgets
- **Entity listing** — filter by domain and state with `ha-ctl entities`
- **Local cache** — automatic caching with configurable TTL
- **JSON output** — structured output for agent consumption (text format also available)

## Installation

```bash
go install github.com/jonson/ha-ctl@latest
```

Or download a binary from [GitHub Releases](https://github.com/jonson/ha-ctl/releases).

## Configuration

Set environment variables:

```bash
export HA_URL=http://homeassistant.local:8123
export HA_TOKEN=your-long-lived-access-token
```

Or create `~/.config/ha-ctl/config.yaml`:

```yaml
ha_url: http://homeassistant.local:8123
ha_token: your-long-lived-access-token
cache_ttl: 300
```

Environment variables take precedence over the config file.

## Usage

```bash
# Find entities by name
ha-ctl find "kitchen"
ha-ctl find "bedroom" --domain light

# Get live entity state
ha-ctl state light.kitchen

# Call a service
ha-ctl call light turn_on --entity light.kitchen
ha-ctl call light turn_on --entity light.kitchen --data '{"brightness":255}'
ha-ctl call climate set_temperature --entity climate.living_room --data '{"temperature":72}'

# Home context summary (optimized for LLM agents)
ha-ctl context
ha-ctl context --full
ha-ctl context --domain sensor,binary_sensor

# List and filter entities
ha-ctl entities --domain light
ha-ctl entities --domain light --state on

# Force cache refresh
ha-ctl cache refresh

# Human-readable output
ha-ctl find "kitchen" --format text
```

## Development

Requires [mise](https://mise.jdx.dev/) for tooling:

```bash
mise install          # install go + goreleaser

mise run build        # build for current platform
mise run test         # run unit tests
mise run test-integration  # run integration tests (requires HA_URL/HA_TOKEN)
mise run snapshot     # build snapshot release with goreleaser
```

Or use make:

```bash
make build            # build for current platform
make build-all        # cross-compile linux amd64 + arm64
make test             # run unit tests
make test-integration # run integration tests
make clean            # remove build artifacts
```

## License

[MIT](LICENSE)
