# go-caffeine

A simple command-line utility to prevent your system from sleeping.

## Features

- Prevents system sleep/screen saver
- Cross-platform support (Windows, macOS, Linux)
- Optional duration timer
- Visual spinner indicator

## Usage

Run indefinitely:
```
go-caffeine
```

Run for specific duration:
```
go-caffeine -t 1h    # Run for 1 hour
go-caffeine -t 30m   # Run for 30 minutes
```

Check version:
```
go-caffeine version
```

## Installation

```
go install github.com/xyztavo/go-caffeine@latest
```

## Building from source

```
git clone https://github.com/xyztavo/go-caffeine.git
cd go-caffeine
go build
```

## License

MIT
