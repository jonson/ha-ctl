---
name: ha-ctl
description: Control Home Assistant devices and entities via the ha-ctl CLI. Use when the user wants to find, inspect, or control smart home devices, check entity states, call Home Assistant services, or get a summary of their home.
compatibility: Requires ha-ctl binary installed and HA_URL/HA_TOKEN configured for a Home Assistant instance.
---

# ha-ctl - Home Assistant CLI

A CLI tool for controlling Home Assistant, designed for use as an LLM agent skill.

## Agent Workflow (read this first)

1. **To find an entity:** `ha-ctl find <name>` — searches entity_id and friendly_name
2. **To control a device:** `ha-ctl call <domain> <service> --entity <entity_id>`
3. **To check live state:** `ha-ctl state <entity_id>`
4. **For broad awareness:** `ha-ctl context` — compact summary (~4KB)
5. **To list/filter:** `ha-ctl entities --domain <d> [--state <s>]`

**DO NOT** use `context` just to find a single entity — use `find` instead.

## Commands

### Find entities by name

```bash
ha-ctl find "kitchen"                     # entities matching "kitchen"
ha-ctl find "bedroom" --domain light      # only lights matching "bedroom"
ha-ctl find "thermostat"                  # find climate entities
```

Case-insensitive substring match against both entity_id and friendly_name. Returns matching entities with state and attributes.

### Get home context

```bash
ha-ctl context                # compact: controllable domains expanded, others as counts
ha-ctl context --full         # all domains expanded (original behavior)
ha-ctl context --domain sensor,binary_sensor  # expand specific domains only
```

Default output shows full details for controllable domains (light, switch, climate, media_player, cover, fan, lock, scene, automation, input_boolean) and entity counts for everything else.

### List entities

```bash
ha-ctl entities                        # all entities
ha-ctl entities --domain light         # only lights
ha-ctl entities --domain light --state on   # lights that are on
ha-ctl entities --state unavailable    # all unavailable entities
ha-ctl entities --refresh              # force cache update first
```

### Get entity state (live)

```bash
ha-ctl state light.kitchen
ha-ctl state climate.living_room
```

Always fetches live from Home Assistant (bypasses cache).

### Call a service

```bash
ha-ctl call light turn_on --entity light.kitchen
ha-ctl call light turn_on --entity light.kitchen --data '{"brightness":255}'
ha-ctl call climate set_temperature --entity climate.living_room --data '{"temperature":72}'
ha-ctl call input_boolean toggle --entity input_boolean.guest_mode
```

### Manage cache

```bash
ha-ctl cache refresh    # force full cache rebuild
```

### Output format

All commands output JSON by default. Use `--format text` for human-readable output.

## Common Service Parameters

| Domain | Service | Parameters |
|--------|---------|------------|
| light | turn_on | brightness (0-255), rgb_color ([r,g,b]), color_temp |
| light | turn_off | - |
| climate | set_temperature | temperature |
| climate | set_hvac_mode | hvac_mode (heat, cool, auto, off) |
| cover | open_cover | - |
| cover | close_cover | - |
| cover | set_cover_position | position (0-100) |
| media_player | media_play | - |
| media_player | media_pause | - |
| media_player | volume_set | volume_level (0.0-1.0) |

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
