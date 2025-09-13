# shenv

`shenv` is a lightweight tool for managing shell environment variables on a per-command. It allows you to define environment variables in a YAML configuration file and load them temporarily when executing specific commands.

## Features

- Temporary environment variable loading
- YAML-based configuration
- No modification of shell profile files required
- Variable expansion support (e.g., `$HOME`, `$PATH`)

## Installation

Install using go install

```bash
$ go install github.com/yourusername/shenv@latest
```

Build from source

```bash
$ go build
```

## Configuration

Create a configuration file at `~/.config/shenv/config.yaml`:

```yaml
env:
  tfenv: |
    TFENV_ROOT="$HOME/.tfenv"
    PATH="$TFENV_ROOT/bin:$PATH"

  nodenv: |
    NODENV_ROOT="$HOME/.nodenv"
    PATH="$NODENV_ROOT/bin:$PATH"
    eval "$(nodenv init -)"
```

## Usage

Basic usage

```bash
$ shenv <environment-name> <command> [args...]
```

Examples

```bash
$ shenv tfenv terraform plan
$ shenv nodenv node --version
```

Using with custom config file

```bash
$ shenv --config /path/to/config.yaml tfenv terraform plan
```

## Examples

### Using with tfenv

Configure tfenv in ~/.config/shenv/config.yaml

```bash
env:
  tfenv: |
    TFENV_ROOT="$HOME/.tfenv"
    PATH="$TFENV_ROOT/bin:$PATH"
```

Use tfenv

```bash
$ shenv tfenv tfenv list
$ shenv tfenv terraform init
```

### Using with custom environment variables

```yaml
env:
  dev: |
    AWS_PROFILE=development
    TF_WORKSPACE=dev
    DEBUG=true
```

## How It Works

1. When you run a command with `shenv`, it:
   - Loads the configuration from your YAML file
   - Sets the specified environment variables
   - Executes your command with the new environment
   - Environment changes only affect the executed command

2. Your shell's environment remains unchanged after the command completes

## License

This project is licensed under the [MIT License](./LICENSE).
