<div align="center">

# 🔧 Foral CLI

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![CI](https://github.com/foral-project/cli/actions/workflows/ci.yml/badge.svg)](https://github.com/foral-project/cli/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8.svg)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/foral-project/cli?include_prereleases)](https://github.com/foral-project/cli/releases)

Governança federada na linha de comando.

[Instalação](#instalação) ·
[Comandos](#comandos) ·
[Protocol](https://github.com/foral-project/protocol) ·
[Governance](https://github.com/foral-project/governance)

</div>

---

## Instalação

### Via Go

```bash
go install github.com/foral-project/cli@latest
```

### Via script (Linux/macOS)

```bash
curl -sfL https://foral-project.github.io/protocol/install.sh | sh
```

### Via GitHub Releases

Baixe o binário para sua plataforma em [Releases](https://github.com/foral-project/cli/releases).

| Plataforma | Binário |
|---|---|
| Linux (x64) | `foral-linux-amd64` |
| Linux (ARM) | `foral-linux-arm64` |
| macOS (x64) | `foral-darwin-amd64` |
| macOS (ARM) | `foral-darwin-arm64` |
| Windows | `foral-windows-amd64.exe` |

## Comandos

### `foral init [path]`

Scaffold de um novo projeto federado:

```bash
foral init my-project --archetype application --owner my-org --lifecycle experimental
```

| Flag | Default | Valores |
|---|---|---|
| `--archetype, -a` | `application` | `application`, `infrastructure`, `bot`, `library`, `service` |
| `--owner, -o` | `foral-project` | RFC 1123 DNS label |
| `--lifecycle, -l` | `experimental` | `experimental`, `production`, `deprecated` |
| `--ci` | `github` | `github`, `gitlab`, `none` |

### `foral validate [file]`

Valida `catalog-info.yaml` contra o [Foral Protocol](https://github.com/foral-project/protocol):

```bash
foral validate                   # valida catalog-info.yaml
foral validate my-catalog.yaml   # valida arquivo específico
foral validate --schema          # apenas schema
foral validate --naming          # apenas naming (RFC 1123)
foral validate --policy          # apenas enums/policies
```

### `foral status`

Dashboard visual de compliance do repositório atual:

```
┌─────────────────────────────────────────────────────┐
│              Foral Compliance Status                │
├─────────────────────────────────────────────────────┤
│  Nome:           my-project                         │
│  Kind:           Component                          │
│  Owner:          my-org                             │
├─────────────────────────────────────────────────────┤
│  ✅  @context presente                              │
│  ✅  apiVersion válido                              │
│  ✅  metadata.name RFC 1123                         │
│  ✅  spec.lifecycle válido                          │
├─────────────────────────────────────────────────────┤
│  Score: [██████████████████████████████] 100% (9/9) │
└─────────────────────────────────────────────────────┘
```

### `foral version`

```bash
foral version          # output legível
foral version --json   # output machine-readable
```

## Como funciona

O CLI consome o [Foral Protocol](https://github.com/foral-project/protocol) offline.
Schemas e policies são embarcados no binário — **zero dependência de rede** para validação local.

Para validação via CI, use os [reusable workflows](https://github.com/foral-project/governance)
do Governance.

## Licença

[Apache-2.0](LICENSE) — SPDX-License-Identifier: Apache-2.0
