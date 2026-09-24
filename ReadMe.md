# Simple CLI Password Manager

Encrypted CLI password manager. Credentials live in `~/.pass/<username>.vault`, sealed with Argon2id + AES-GCM.

## Features

- Master-password protected vault (one file per user)
- Add / list / get / delete entries
- Password generation (`pass generate` or from the vault menu)
- Hidden master-password input

## Build

```bash
make build
# or
go build -o .bin/cmd
```

## Usage

Interactive vault:

```bash
make start-app
# or
.bin/cmd start
```

Options: create account, login, then use the vault menu (`add`, `list`, `get`, `delete`, `generate`, `exit`).

Generate a password without unlocking:

```bash
.bin/cmd generate -l 20 -d -s
```

## Security notes

- Choose a strong master password; it cannot be recovered.
- Vault files are mode `0600` under `~/.pass/`.
- The master password is derived with Argon2id; entries are encrypted with AES-256-GCM.
- This is a personal learning project — not a replacement for a audited password manager.
