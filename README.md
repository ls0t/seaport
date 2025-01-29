# seaport

![seaport definition](docs/seaport-definition.png)

Seaport sets up port forwarding to allow directly reaching services, whether that be torrent clients, web servers, or games. In many cases, port forwarding will make connectivity possible where it was not before, leading to better utilization and lower latency. If the port changes, seaport will update your client and optionally notify you.

Seaport is modular. If you know of an environment that seaport should support, please file a ticket or submit a PR.

## Quick Start

First create your config file and save it as *seaport.yaml*. Then copy/paste one config block example below.

**For protonvpn with qbittorent**

```yaml
source:
  name: "protonvpn"

actions:
  - name: qbittorrent
    options:
      # adjust the options here if you have changed the qbittorrent defaults
      host: http://localhost:8080
      username: admin
      password: adminadmin
```

**For getting a port from your router (or generic vpn) with qbittorent**

```yaml
source:
  name: "natpmp"
  options:
    gatewayIP: <insert gateway>

actions:
  - name: qbittorrent
    options:
      # adjust the options here if you have changed the qbittorrent defaults
      host: http://localhost:8080
      username: admin
      password: adminadmin
```

## Releases

Download the latest release from github releases.

## Key Concepts

_Sources_ are where IP+port comes from for port forwarding. Usually from your infrastructure, like your router or VPN provider.

_Actions_ are plugins that automatically configure external clients using programmatic means, like your torrent client. You need at least one action for seaport to be useful to you.

_Notifiers_ are optional ways to send human-readable updates.

You can only have one source, but many actions or notifiers.

## Full Config Example

See [seaport-example.yaml](seaport-example.yaml).

## Supported Plugins

### Sources

#### protonvpn

No configuration needed.

```yaml
source:
  name: protonvpn
```

#### natpmp

```yaml
source:
  name: natpmp
  options:
    gatewayIP: <insert IP>
```

### Actions

#### qbittorrent

```yaml
actions:
  - name: qbittorrent
    options:
      host: http://localhost:8080
      username: admin
      password: adminadmin
```

### Notifiers

#### discord

```yaml
notifiers:
  - name: discord
    options:
      webhook: <insert your webhook URL here like https://discord.com/api/webhooks/...>
```

Alternatively you can specify the webhook by setting the environment variable `SEAPORT_WEBHOOK_URL`.