# Prism 💍
> A TUI for reviewing GitHub Pull Requests in your terminal

## Features
- Review PRs without leaving the terminal
- Vim-like keybindings for fast navigation
- Diff view, commenting, and merge support
- CI status checks
- Auto-save review progress

## Installation
```bash
go install github.com/urabexon/prism@latest
```

## Usage
1. Run in current repository
```bash
prism
```
2. Specify repository
```bash
prism owner/repo
```

## Keybindings
### PR List
| Key     | Action       |
| ------- | ------------ |
| `j/k`   | Navigate     |
| `Enter` | Open PR      |
| `m`     | Merge        |
| `c`     | CI checks    |
| `q`     | Quit         |

### Diff View
| Key     | Action          |
| ------- | --------------- |
| `j/k`   | Scroll          |
| `n/N`   | Next/Prev hunk  |
| `V`     | Visual select   |
| `space` | Mark reviewed   |
| `esc`   | Back            |

## Requirements
- Go 1.21+
- GitHub CLI (gh) authenticated

## License
MIT
