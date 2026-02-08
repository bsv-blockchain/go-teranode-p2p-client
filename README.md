<div align="center">

# 🛰&nbsp;&nbsp;go-teranode-p2p-client

**Go library for connecting to Teranode's P2P gossip network as a client/subscriber.**

<br/>

<a href="https://github.com/bsv-blockchain/go-teranode-p2p-client/releases"><img src="https://img.shields.io/github/release-pre/bsv-blockchain/go-teranode-p2p-client?include_prereleases&style=flat-square&logo=github&color=black" alt="Release"></a>
<a href="https://golang.org/"><img src="https://img.shields.io/github/go-mod/go-version/bsv-blockchain/go-teranode-p2p-client?style=flat-square&logo=go&color=00ADD8" alt="Go Version"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/license-OpenBSV-blue?style=flat-square&logo=springsecurity&logoColor=white" alt="License"></a>

<br/>

<table align="center" border="0">
  <tr>
    <td align="right">
       <code>CI / CD</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/go-teranode-p2p-client/actions"><img src="https://img.shields.io/github/actions/workflow/status/bsv-blockchain/go-teranode-p2p-client/fortress.yml?branch=main&label=build&logo=github&style=flat-square" alt="Build"></a>
       <a href="https://github.com/bsv-blockchain/go-teranode-p2p-client/actions"><img src="https://img.shields.io/github/last-commit/bsv-blockchain/go-teranode-p2p-client?style=flat-square&logo=git&logoColor=white&label=last%20update" alt="Last Commit"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Quality</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://goreportcard.com/report/github.com/bsv-blockchain/go-teranode-p2p-client"><img src="https://goreportcard.com/badge/github.com/bsv-blockchain/go-teranode-p2p-client?style=flat-square" alt="Go Report"></a>
       <a href="https://codecov.io/gh/bsv-blockchain/go-teranode-p2p-client"><img src="https://codecov.io/gh/bsv-blockchain/go-teranode-p2p-client/branch/main/graph/badge.svg?style=flat-square" alt="Coverage"></a>
    </td>
  </tr>

  <tr>
    <td align="right">
       <code>Security</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://scorecard.dev/viewer/?uri=github.com/bsv-blockchain/go-teranode-p2p-client"><img src="https://api.scorecard.dev/projects/github.com/bsv-blockchain/go-teranode-p2p-client/badge?style=flat-square" alt="Scorecard"></a>
       <a href=".github/SECURITY.md"><img src="https://img.shields.io/badge/policy-active-success?style=flat-square&logo=security&logoColor=white" alt="Security"></a>
    </td>
    <td align="right">
       &nbsp;&nbsp;&nbsp;&nbsp; <code>Community</code> &nbsp;&nbsp;
    </td>
    <td align="left">
       <a href="https://github.com/bsv-blockchain/go-teranode-p2p-client/graphs/contributors"><img src="https://img.shields.io/github/contributors/bsv-blockchain/go-teranode-p2p-client?style=flat-square&color=orange" alt="Contributors"></a>
       <a href="https://github.com/sponsors/bsv-blockchain"><img src="https://img.shields.io/badge/sponsor-BSV-181717.svg?logo=github&style=flat-square" alt="Sponsor"></a>
    </td>
  </tr>
</table>

</div>

<br/>
<br/>

<div align="center">

### <code>Project Navigation</code>

</div>

<table align="center">
  <tr>
    <td align="center" width="33%">
       📦&nbsp;<a href="#-installation"><code>Installation</code></a>
    </td>
    <td align="center" width="33%">
       🧪&nbsp;<a href="#-examples--tests"><code>Examples&nbsp;&&nbsp;Tests</code></a>
    </td>
    <td align="center" width="33%">
       📚&nbsp;<a href="#-documentation"><code>Documentation</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🛠️&nbsp;<a href="#-code-standards"><code>Code&nbsp;Standards</code></a>
    </td>
    <td align="center">
       ⚡&nbsp;<a href="#-benchmarks"><code>Benchmarks</code></a>
    </td>
    <td align="center">
       🤖&nbsp;<a href="#-ai-usage--assistant-guidelines"><code>AI&nbsp;Usage</code></a>
    </td>
  </tr>
  <tr>
    <td align="center">
       🤝&nbsp;<a href="#-contributing"><code>Contributing</code></a>
    </td>
    <td align="center">
       👥&nbsp;<a href="#-maintainers"><code>Maintainers</code></a>
    </td>
    <td align="center">
       ⚖️&nbsp;<a href="#-license"><code>License</code></a>
    </td>
  </tr>
</table>
<br/>

## 📦 Installation

**go-teranode-p2p-client** requires a [supported release of Go](https://golang.org/doc/devel/release.html#policy).
```shell script
go get -u github.com/bsv-blockchain/go-teranode-p2p-client
```

<br/>

## 📚 Documentation
- **API Reference** – Dive into the godocs at [pkg.go.dev/github.com/bsv-blockchain/go-teranode-p2p-client](https://pkg.go.dev/github.com/bsv-blockchain/go-teranode-p2p-client)

## Features
- Re-exports canonical message types from Teranode's P2P package
- Embedded bootstrap peers for mainnet, testnet, and STN
- Persistent P2P identity (key management)
- Topic name helpers

## Message Types

- `BlockMessage` - New block announcements
- `SubtreeMessage` - Transaction batch (subtree) announcements
- `RejectedTxMessage` - Rejected transaction notifications
- `NodeStatusMessage` - Node status updates

## Topics

| Constant          | Topic Name    |
|-------------------|---------------|
| `TopicBlock`      | `block`       |
| `TopicSubtree`    | `subtree`     |
| `TopicRejectedTx` | `rejected-tx` |
| `TopicNodeStatus` | `node_status` |

Use `TopicName(network, topic)` to construct full topic names (e.g., `teranode/bitcoin/1.0.0/mainnet-block`).

## Networks

| Constant             | Network       |
|----------------------|---------------|
| `NetworkMainnet`     | `mainnet`     |
| `NetworkTestnet`     | `testnet`     |
| `NetworkSTN`         | `stn`         |
| `NetworkTeratestnet` | `teratestnet` |

> **Note:** `NetworkTeratestnet` requires manual bootstrap peer configuration via `Config.MsgBus.BootstrapPeers`.

## Bootstrap Peers

Embedded bootstrap peers for `main`, `test`, and `stn` networks are automatically applied when using `Config.Initialize()`.

## Usage

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    p2pclient "github.com/bsv-blockchain/go-teranode-p2p-client"
)

func main() {
    // Create a context that cancels on interrupt signals
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    // Configure the P2P client
    cfg := p2pclient.Config{
        Network:     "main",   // Options: "main", "test", "stn"
        StoragePath: "./data", // Persistent storage for keys and peer cache
    }

    // Initialize the client (loads/generates P2P identity, connects to bootstrap peers)
    client, err := cfg.Initialize(ctx, "my-app")
    if err != nil {
        slog.Error("failed to initialize P2P client", slog.String("error", err.Error()))
        os.Exit(1)
    }
    defer client.Close()

    slog.Info("connected to Teranode P2P network",
        slog.String("peer_id", client.GetID()),
        slog.String("network", client.GetNetwork()),
    )

    // Subscribe to block announcements (returns typed channel)
    blocks := client.SubscribeBlocks(ctx)

    for block := range blocks {
        slog.Info("new block",
            slog.Uint64("height", uint64(block.Height)),
            slog.String("hash", block.Hash),
        )
    }
}
```

<br/>

<details>
<summary><strong><code>Development Build Commands</code></strong></summary>
<br/>

Get the [MAGE-X](https://github.com/mrz1836/mage-x) build tool for development:
```shell script
go install github.com/mrz1836/mage-x/cmd/magex@latest
```

View all build commands

```bash script
magex help
```

</details>

<details>
<summary><strong>Repository Features</strong></summary>
<br/>

This repository includes 25+ built-in features covering CI/CD, security, code quality, developer experience, and community tooling.

**[View the full Repository Features list →](.github/docs/repository-features.md)**

</details>

<details>
<summary><strong><code>Library Deployment</code></strong></summary>
<br/>

This project uses [goreleaser](https://github.com/goreleaser/goreleaser) for streamlined binary and library deployment to GitHub. To get started, install it via:

```bash
brew install goreleaser
```

The release process is defined in the [.goreleaser.yml](.goreleaser.yml) configuration file.


Then create and push a new Git tag using:

```bash
magex version:bump push=true bump=patch branch=main
```

This process ensures consistent, repeatable releases with properly versioned artifacts and citation metadata.

</details>

<details>
<summary><strong><code>Pre-commit Hooks</code></strong></summary>
<br/>

Set up the Go-Pre-commit System to run the same formatting, linting, and tests defined in [AGENTS.md](.github/AGENTS.md) before every commit:

```bash
go install github.com/mrz1836/go-pre-commit/cmd/go-pre-commit@latest
go-pre-commit install
```

The system is configured via modular env files in [`.github/env/`](.github/env/README.md) and provides 17x faster execution than traditional Python-based pre-commit hooks. See the [complete documentation](http://github.com/mrz1836/go-pre-commit) for details.

</details>

<details>
<summary><strong>GitHub Workflows</strong></summary>
<br/>

All workflows are driven by modular configuration in [`.github/env/`](.github/env/README.md) — no YAML editing required.

**[View all workflows and the control center →](.github/docs/workflows.md)**

</details>

<details>
<summary><strong><code>Updating Dependencies</code></strong></summary>
<br/>

To update all dependencies (Go modules, linters, and related tools), run:

```bash
magex deps:update
```

This command ensures all dependencies are brought up to date in a single step, including Go modules and any tools managed by [MAGE-X](https://github.com/mrz1836/mage-x). It is the recommended way to keep your development environment and CI in sync with the latest versions.

</details>

<br/>

## 🧪 Examples & Tests

All unit tests and examples run via [GitHub Actions](https://github.com/bsv-blockchain/go-teranode-p2p-client/actions) and use [Go version 1.25.x](https://go.dev/doc/go1.25). View the [configuration file](.github/workflows/fortress.yml).

Run all tests (fast):

```bash script
magex test
```

Run all tests with race detector (slower):
```bash script
magex test:race
```

<br/>

## ⚡ Benchmarks

Run the Go benchmarks:

```bash script
magex bench
```

> **Note:** Comprehensive benchmarks for P2P operations (peer discovery, message throughput, connection establishment) are planned for future releases. The current focus is on correctness and stability of the networking implementation.

<br/>

## 🛠️ Code Standards
Read more about this Go project's [code standards](.github/CODE_STANDARDS.md).

<br/>

## 🤖 AI Usage & Assistant Guidelines
Read the [AI Usage & Assistant Guidelines](.github/tech-conventions/ai-compliance.md) for details on how AI is used in this project and how to interact with AI assistants.

<br/>

## 👥 Maintainers
| [<img src="https://github.com/icellan.png" height="50" alt="Siggi" />](https://github.com/icellan) | [<img src="https://github.com/galt-tr.png" height="50" alt="Galt" />](https://github.com/galt-tr) | [<img src="https://github.com/mrz1836.png" height="50" alt="MrZ" />](https://github.com/mrz1836) |
|:--------------------------------------------------------------------------------------------------:|:-------------------------------------------------------------------------------------------------:|:------------------------------------------------------------------------------------------------:|
|                                [Siggi](https://github.com/icellan)                                 |                                [Dylan](https://github.com/galt-tr)                                |                                [MrZ](https://github.com/mrz1836)                                 |

<br/>

## 🤝 Contributing
View the [contributing guidelines](.github/CONTRIBUTING.md) and please follow the [code of conduct](.github/CODE_OF_CONDUCT.md).

### How can I help?
All kinds of contributions are welcome :raised_hands:!
The most basic way to show your support is to star :star2: the project, or to raise issues :speech_balloon:.

[![Stars](https://img.shields.io/github/stars/bsv-blockchain/go-teranode-p2p-client?label=Please%20like%20us&style=social&v=1)](https://github.com/bsv-blockchain/go-teranode-p2p-client/stargazers)

<br/>

## 📝 License

[![License](https://img.shields.io/badge/license-OpenBSV-blue?style=flat&logo=springsecurity&logoColor=white)](LICENSE)
