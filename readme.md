# GitUp

**GitUp** is a simple, user-friendly Git client written in Go.

The goal of GitUp is to make common Git operations easier and more approachable while providing a clean command-line interface.

> 🚧 **GitUp is currently in early development.**
> Expect missing features, changes, and the occasional spectacular bug.

## Features

Currently implemented:

* `init` — Initialize a Git repository
* `commit` — Commit specific files, or all changes
* `commit --dynamic_file` — Interactively pick which files go into a commit
* `login` / `logout` / `whoami` — GitHub authentication via OAuth device flow
* Interactive confirmation dialogs
* Colored/status-style logging
* Cross-platform Go foundation
* Default `.gitignore` generation *(in development)*

Planned commands include:

```text
gitup init
gitup status
gitup commit
gitup clone
gitup publish
gitup put aside
gitup get aside
gitup get
gitup sync
gitup erase history <commit>
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
│   ├── auth/
│   │   ├── auth.go        # Device flow client + credential storage
│   │   └── commands.go    # login / logout / whoami
│   ├── git/
│   │   ├── init.go
│   │   ├── commit.go
│   │   ├── select.go       # Interactive --dynamic_file picker
│   │   ├── gitignore.go
│   │   └── misc.go
│   └── helper/
│       ├── log.go
│       ├── diag.go
│       ├── browser.go
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

Commit all changes:

```bash
gitup commit "Fix filesystem"
```

Commit specific files:

```bash
gitup commit src/foo.c src/bar.c "Fix filesystem"
```

Pick which files to commit interactively:

```bash
gitup commit --dynamic_file "Fix filesystem"
```

This opens a checklist of every changed file. Type numbers to toggle files
(`1 3 4`), `a` to select all, `n` to select none, `c` to confirm and commit,
or `q` to cancel.

Log in to GitHub:

```bash
gitup login
```

This opens your browser and walks you through GitHub's device authorization
flow — no password is ever entered into GitUp.

Check who you're logged in as:

```bash
gitup whoami
```

Log out:

```bash
gitup logout
```

Display help:

```bash
gitup --help
```

## Authentication

GitUp authenticates with GitHub through a **GitHub App using the OAuth
device flow**. Running `gitup login` will:

1. Ask GitHub for a short one-time code.
2. Open your browser (or print a URL + code if it can't) to
   `https://github.com/login/device`.
3. Wait for you to enter the code and approve access.
4. Store the resulting access + refresh token locally, associated with your
   GitHub username.

GitUp never asks for or stores your GitHub password. Access tokens are
short-lived and refreshed automatically in the background — if a token is
ever lost or leaked, it stops working on its own.

Credentials are currently stored in a local config file
(`~/.config/gitup/credentials.json` on Linux, the equivalent
`AppData`/`Application Support` location on Windows/macOS). Moving this to
OS-native secure credential storage (Keychain / Credential Manager / Secret
Service) is a planned follow-up.

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