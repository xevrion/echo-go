# echo-go  
**A simple peer‑to‑peer (P2P) terminal UI chat application written in Go**  

![Build Status](https://img.shields.io/github/actions/workflow/status/xevrion/echo-go/ci.yml?branch=main&label=build) ![Coverage](https://img.shields.io/codecov/c/github/xevrion/echo-go?label=coverage) ![Version](https://img.shields.io/github/v/tag/xevrion/echo-go?label=version) ![License](https://img.shields.io/github/license/xevrion/echo-go)  

[Demo ▶️](#demo) • [Documentation 📚](#documentation) • [Issues 🐞](https://github.com/xevrion/echo-go/issues)  

---  

## Overview  

`echo-go` is a lightweight, terminal‑based chat client that connects directly to other peers without a central server. It is built with Go’s standard library and a minimal set of third‑party dependencies, making it fast, portable, and easy to extend.  

* **Zero‑install** – run a single binary on any platform that supports Go.  
* **True P2P** – messages travel directly between peers; no third‑party server is required.  
* **TUI‑first** – a clean, ncurses‑style interface works over SSH, Docker, or any terminal.  

Target audience: developers, sysadmins, and hobbyists who need quick, secure, ad‑hoc chat sessions on the command line.

Current stable version: **v0.1.0** (2026‑01‑27)  

---  

## Features  

| Feature | Description | Status |
|---------|-------------|--------|
| **Peer discovery** | Connect to peers via IP:port or via a simple “bootstrap” file. | ✅ Stable |
| **End‑to‑end encryption** | Optional NaCl box encryption for all traffic. | ✅ Stable |
| **Terminal UI** | Full‑screen chat view with message history, timestamps, and user list. | 🚧 Beta |
| **File transfer** | Drag‑and‑drop (via command) to send binary files. | 🚧 Experimental |
| **Multi‑room support** | Create isolated chat rooms on the same network. | ❌ Planned |
| **Message persistence** | Local SQLite log of all received/sent messages. | ❌ Planned |
| **Cross‑platform** | Works on Linux, macOS, and Windows (via WSL or native). | ✅ Stable |

---  

## Tech Stack  

| Layer | Technology | Reason |
|-------|------------|--------|
| **Language** | Go 1.22+ | Strong concurrency, static binary output |
| **UI** | [tview](https://github.com/rivo/tview) + [tcell](https://github.com/gdamore/tcell) | Rich terminal widgets, cross‑platform |
| **Networking** | Go `net` package + optional [libsodium-go](https://github.com/GoKillers/libsodium-go) for encryption | Minimal dependencies, full control |
| **Persistence** | SQLite (via `github.com/mattn/go-sqlite3`) | Lightweight local storage |
| **Testing** | Go `testing` + `github.com/stretchr/testify` | Familiar, expressive assertions |
| **CI/CD** | GitHub Actions | Automated build, test, and coverage reports |

---  

## Architecture  

```
echo-go/
├── cmd/                # main entry point (CLI)
│   └── echo-go.go
├── internal/
│   ├── ui/             # TUI components (tview primitives)
│   ├── net/            # Peer connection handling, discovery
│   ├── crypto/         # Encryption utilities
│   └── storage/        # SQLite logger
├── pkg/                # Public reusable packages (e.g., message structs)
├── config/             # Default config files & templates
└── scripts/            # Helper scripts (dev, lint, release)
```

* **CLI (`cmd/echo-go.go`)** parses flags, loads configuration, and starts the UI.  
* **`internal/net`** opens a listening socket, dials peers, and multiplexes streams.  
* **`internal/crypto`** wraps NaCl box for key generation, encryption, and decryption.  
* **`internal/ui`** builds the terminal interface, handles user input, and renders messages.  
* **`internal/storage`** persists chat history in a local SQLite DB (`$HOME/.echo-go/history.db`).  

Data flow: UI → Message struct → (optional) encrypt → network layer → remote peer → (optional) decrypt → UI.  

---  

## Getting Started  

### Prerequisites  

| Tool | Minimum version |
|------|-----------------|
| Go | 1.22 |
| Git | any |
| Make (optional) | 4.0+ |
| A terminal that supports true‑color (most modern terminals) | — |

### Installation  

#### 1. Install via `go install` (recommended)  

```bash
go install github.com/xevrion/echo-go@latest
# Binary will be placed in $(go env GOPATH)/bin
```

#### 2. Build from source  

```bash
git clone https://github.com/xevrion/echo-go.git
cd echo-go
make build          # uses the provided Makefile
./bin/echo-go --help
```

#### 3. Docker (experimental)  

```bash
docker pull ghcr.io/xevrion/echo-go:latest
docker run -it --rm \
  -v $HOME/.echo-go:/root/.echo-go \
  ghcr.io/xevrion/echo-go:latest --listen :9000
```

### Configuration  

`echo-go` reads configuration from `$HOME/.echo-go/config.yaml` (created on first run) and environment variables.  

#### Environment variables  

| Variable | Description | Default |
|----------|-------------|---------|
| `ECHO_PORT` | Local listening port | `9000` |
| `ECHO_BOOTSTRAP` | Path to a file containing known peer addresses (one per line) | `$HOME/.echo-go/bootstrap.txt` |
| `ECHO_ENCRYPT` | Enable end‑to‑end encryption (`true`/`false`) | `true` |
| `ECHO_LOG_LEVEL` | Log verbosity (`debug`, `info`, `warn`, `error`) | `info` |

#### Example `.env`  

```dotenv
ECHO_PORT=9000
ECHO_BOOTSTRAP=$HOME/.echo-go/bootstrap.txt
ECHO_ENCRYPT=true
ECHO_LOG_LEVEL=info
```

Copy the above into a file named `.env` in the project root and run:

```bash
export $(cat .env | xargs)   # loads variables into the shell
```

---  

## Usage  

```bash
# Start the application (will listen on $ECHO_PORT)
echo-go
```

### Command‑line flags  

| Flag | Description | Example |
|------|-------------|---------|
| `--listen <addr>` | Override listening address (default: `:9000`) | `--listen :8000` |
| `--peer <addr>` | Connect to a remote peer immediately | `--peer 192.168.1.42:9000` |
| `--no-encrypt` | Disable encryption (useful for debugging) | `--no-encrypt` |
| `--config <path>` | Path to a custom YAML config file | `--config ./myconfig.yaml` |
| `--help` | Show help message | — |

### Typical workflow  

1. **Start your instance**  

   ```bash
   echo-go
   ```

2. **Add peers** – either edit `$HOME/.echo-go/bootstrap.txt` with `IP:PORT` lines and press **Ctrl+R** in the UI, or use the `--peer` flag on launch.  

3. **Chat** – type your message and press **Enter**. Messages appear with timestamps and the sender’s nickname (default is the OS username).  

4. **File transfer** (experimental)  

   ```bash
   /send /path/to/file.txt
   ```

   The recipient will see a prompt to accept or reject.  

### Screenshots  

<div align="center">
  <img src="https://raw.githubusercontent.com/xevrion/echo-go/main/assets/screenshot.png" alt="Echo-go TUI" width="800"/>
</div>

---  

## Development  

### Clone the repository  

```bash
git clone https://github.com/xevrion/echo-go.git
cd echo-go
```

### Run tests  

```bash
make test          # runs `go test ./...` with coverage
```

### Lint & format  

```bash
make lint          # uses golangci-lint
make fmt           # runs `go fmt ./...`
```

### Debugging  

* Enable verbose logs: `ECHO_LOG_LEVEL=debug echo-go`  
* Use `pprof` by starting the binary with `--pprof :6060` (future flag).  

### Contributing  

1. Fork the repo and create a feature branch (`git checkout -b feat/your-feature`).  
2. Write tests for any new functionality.  
3. Ensure `make lint && make test` passes.  
4. Open a Pull Request with a clear description of the change.  

Please follow the **Go Code Review Comments** style guide and keep commit messages concise.  

---  

## Deployment  

### Production binary  

```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=v0.1.0" -o echo-go ./cmd
```

Distribute the binary via your preferred channel (GitHub Releases, internal artifact repo, etc.).  

### Docker (recommended for CI)  

```Dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -ldflags="-s -w -X main.version=v0.1.0" -o /echo-go ./cmd

FROM alpine:3.19
COPY --from=builder /echo-go /usr/local/bin/echo-go
ENTRYPOINT ["echo-go"]
```

Build and push:

```bash
docker build -t ghcr.io/xevrion/echo-go:0.1.0 .
docker push ghcr.io/xevrion/echo-go:0.1.0
```

### Performance tips  

* Run with `GOMAXPROCS=$(nproc)` on multi‑core machines.  
* Keep the message history size reasonable (`sqlite3` pragma `journal_mode=WAL`).  

---  

## API Documentation  

`echo-go` does not expose a public HTTP/REST API. All communication occurs over encrypted TCP streams using a custom binary protocol defined in `internal/net/protocol.go`.  

Key protocol messages (for developers extending the project):  

| Type | Direction | Payload | Description |
|------|-----------|---------|-------------|
| `HELLO` | Outbound | `{Version:string, Nick:string, PubKey:[]byte}` | Peer handshake |
| `MSG` | Bidirectional | `{Timestamp:int64, From:string, Body:string}` | Chat message |
| `FILE_OFFER` | Outbound | `{FileName:string, Size:int64, Hash:[]byte}` | Initiate file transfer |
| `FILE_ACK` | Inbound | `{Accepted:bool}` | Receiver response |
| `PING` / `PONG` | Keep‑alive | — | Detect dead connections |

For deeper details, see the source file `internal/net/protocol.go`.  

---  

## Roadmap  

| Milestone | Target date | Planned features |
|-----------|-------------|------------------|
| **v0.2.0** | 2026‑04‑01 | Multi‑room support, message persistence, UI polish |
| **v0.3.0** | 2026‑07‑15 | Web‑socket bridge, mobile client, plugin system |
| **v1.0.0** | 2027‑01‑01 | Full test coverage (>90 %), stable public API, official Docker images |

---  

## License & Credits  

**MIT License** – see the [LICENSE](https://github.com/xevrion/echo-go/blob/main/LICENSE) file for details.  

### Contributors  

- **xevrion** – project creator & lead maintainer  
- *(Add your name here when you contribute!)*  

### Acknowledgments  

- The **tview** and **tcell** projects for making terminal UIs a joy.  
- **NaCl** for simple, battle‑tested encryption primitives.  
- All open‑source libraries listed in `go.mod`.  

---  

*Happy chatting! 🎉*  