# GitUp

**GitUp** is a simple, user-friendly Git client written in Go.

The goal of GitUp is to make common Git operations easier and more approachable while providing a clean command-line interface.

> 🚧 **GitUp is currently in early development.**
> Expect missing features, changes, and the occasional spectacular bug.

## Features

Currently implemented:

* `init` — Initialize a Git repository
* Interactive confirmation dialogs
* Colored/status-style logging
* Cross-platform Go foundation
* Default `.gitignore` generation *(in development)*

Planned commands include:

```text
gitup init
gitup status
gitup add
gitup commit
gitup clone
gitup publish
gitup put aside
gitup get aside
```

## Why GitUp?

Git is incredibly powerful, but its command-line interface can be intimidating when you're starting out.

GitUp aims to provide a simpler interface for common Git workflows without hiding what is actually happening.

## Built With

* **Go** — Main programming language
* **go-git** — Git implementation/library for Go

## Project Structure

```text
GitUp/
├── src/
│   ├── main.go
│   ├── git/
│   │   └── init.go
│   └── helper/
│       ├── log.go
│       └── ...
├── go.mod
├── go.sum
└── README.md
```

## Building

Make sure you have Go installed.

Clone the repository and build it:

```bash
git clone <repository-url>
cd GitUp
go build ./src
```

You can also run GitUp directly during development:

```bash
go run ./src
```

## Usage

Initialize a repository:

```bash
gitup init ./my-project
```

Display help:

```bash
gitup --help
```

## Cross-Platform

GitUp is written in Go with cross-platform support in mind.

Builds can be produced for different operating systems and architectures using Go's cross-compilation support.

## Development

GitUp is being developed as a learning project as well as a useful tool.

The project is intentionally built from the ground up around Go while using `go-git` for the underlying Git functionality.

Contributions, ideas, bug reports, and feedback are welcome once the project is ready for them.

## License

MIT License.
See LICENSE.txt for more details.

---

**GitUp — Git, but going up. 🚀**

<p align="center">by: New Technologies Software (NtinosTheGamer2324)</p>