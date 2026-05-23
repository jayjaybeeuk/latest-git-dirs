# latestdirs

A fast, concurrent directory scanner that finds the most recently modified project directories. Written in Go for single-binary deployment.

## Features

- Scans immediate child directories concurrently
- Determines latest modification time via filesystem walk
- Sorts newest first
- Configurable worker pool
- Excludes common noise directories (`.git`, `node_modules`, `bin`, `obj`)
- Supports JSON output
- Compiles to a single binary

## Usage

```bash
# Scan current directory
latestdirs

# Scan a specific path
latestdirs /path/to/projects

# Top 10, JSON output
latestdirs -top 10 -json /path/to/projects

# Custom depth and workers
latestdirs -max-depth 5 -workers 8 .
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-top` | 20 | Number of results to show |
| `-max-depth` | 10 | Maximum directory depth to scan |
| `-json` | false | Output as JSON |
| `-workers` | CPU count | Number of concurrent workers |

## Build

```bash
# Local
go build ./cmd/latestdirs

# Cross-platform
GOOS=windows GOARCH=amd64 go build -o dist/latestdirs.exe ./cmd/latestdirs
GOOS=linux GOARCH=amd64 go build -o dist/latestdirs-linux ./cmd/latestdirs
GOOS=darwin GOARCH=arm64 go build -o dist/latestdirs-macos ./cmd/latestdirs
```

## Project Structure

```
latestdirs/
├── cmd/latestdirs/main.go      # CLI entry point
├── internal/
│   ├── config/config.go        # Configuration struct
│   ├── git/git.go              # Git mode (planned)
│   ├── model/result.go         # Result type
│   ├── output/output.go        # Table and JSON formatters
│   └── scanner/scanner.go      # Filesystem walker
├── go.mod
└── README.md
```

## Performance

This design achieves speed through:

- Concurrent scanning of top-level projects via worker pool
- Early skipping of excluded directories
- Zero external dependencies
- Go's efficient `filepath.WalkDir`
- Bounded goroutine creation

### Profiling

```bash
go test ./... -bench=. -benchmem
go test -cpuprofile cpu.prof
go tool pprof cpu.prof
```

## Roadmap

### Milestone 1 ✅
- Basic filesystem scanner
- Concurrent worker pool
- Table output

### Milestone 2
- Git mode (latest commit time)
- JSON schema
- Configurable exclusions
- Depth limiting

### Milestone 3
- Benchmark suite
- Compare against .NET implementation
- Optimise allocations

### Milestone 4
- GitHub Actions CI
- GoReleaser
- Published binaries

### Future
- Fuzzy search
- Watch mode
- Interactive TUI
- Shell completions
- Caching/index mode
- Sort by branch / repo size
- Coloured output

## License

MIT
