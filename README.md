# Foral CLI

Governança federada na linha de comando.

```bash
foral init my-project
foral validate
foral status
```

## Instalação

### Via script (Linux/macOS)

```bash
curl -sfL https://foral-project.github.io/protocol/install.sh | sh
```

### Via Go

```bash
go install github.com/foral-project/cli@latest
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

Flags:
- `--archetype, -a` — Tipo do projeto: `application`, `infrastructure`, `bot`, `library`, `service`
- `--owner, -o` — Owner do projeto (RFC 1123 DNS label)
- `--lifecycle, -l` — Ciclo de vida: `experimental`, `production`, `deprecated`
- `--ci` — Plataforma CI: `github`, `gitlab`, `none`

### `foral validate [file]`

Valida `catalog-info.yaml` contra o Foral Protocol:

```bash
foral validate                   # valida catalog-info.yaml
foral validate my-catalog.yaml   # valida arquivo específico
foral validate --schema          # apenas validação de schema
foral validate --naming          # apenas validação de naming
```

### `foral status`

Dashboard de compliance do repositório atual:

```
┌─────────────────────────────────────────────────────┐
│              Foral Compliance Status                │
├─────────────────────────────────────────────────────┤
│  Nome:           my-project                          │
│  Kind:           Component                           │
│  Owner:          my-org                              │
├─────────────────────────────────────────────────────┤
│  ✅  @context presente                             │
│  ✅  apiVersion válido                             │
│  ✅  metadata.name RFC 1123                        │
├─────────────────────────────────────────────────────┤
│  Score: [██████████████████████████████] 100% (9/9) │
└─────────────────────────────────────────────────────┘
```

### `foral version`

```bash
foral version          # output legível
foral version --json   # output machine-readable
```

## Protocol

Este CLI consome o [Foral Protocol](https://github.com/foral-project/protocol)
e os [reusable workflows](https://github.com/foral-project/governance) da governance.

Schemas e policies são baixados via HTTP do GitHub Pages — zero dependência local.

## Licença

Apache-2.0
