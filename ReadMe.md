# Simple CLI Password Manager

Encrypted CLI password manager with a [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI. Credentials live in `~/.pass/<username>.vault`, sealed with Argon2id + AES-GCM.

## Features

- Interactive TUI (`pass start`): login, create account, browse vault
- Master-password protected vault (one file per user)
- Add / list / view / delete entries
- Password generation (`pass generate` or `g` in the TUI)
- Hidden master-password input

## Build

```bash
make build
# or
go build -o .bin/cmd .
```

## Usage

```bash
make start-app
# or
.bin/cmd start
```

**TUI keys**

| Screen | Keys |
|--------|------|
| Welcome | ↑/↓, enter, q |
| Auth | tab, enter, esc |
| Vault | ↑/↓, `/` filter, enter view, a add, d delete, g generate, q quit |
| Add | tab, enter save, esc |
| Detail | esc back, d delete |

Generate without the TUI:

```bash
.bin/cmd generate -l 20 -d -s
```

## Security notes

- Choose a strong master password; it cannot be recovered.
- Vault files are mode `0600` under `~/.pass/`.
- The master password is derived with Argon2id; entries are encrypted with AES-256-GCM.
- This is a personal learning project — not a replacement for an audited password manager.
