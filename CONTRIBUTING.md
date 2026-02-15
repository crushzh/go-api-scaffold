# Contributing

Thank you for your interest in contributing! Here's how to get started.

## Reporting Bugs

- Search [existing issues](https://github.com/mrzhoong/go-api-scaffold/issues) first
- If not found, open a new issue with:
  - Go version (`go version`)
  - OS and architecture
  - Steps to reproduce
  - Expected vs actual behavior

## Suggesting Features

Open an issue with the `feature` label describing:
- The problem you're trying to solve
- Your proposed solution
- Any alternatives you've considered

## Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit with a clear message: `git commit -m "feat: add my feature"`
7. Push and open a PR

### Commit Convention

```
feat:     New feature
fix:      Bug fix
refactor: Code refactoring
docs:     Documentation
style:    Code formatting
test:     Tests
chore:    Build / tooling
```

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Add comments for exported functions
- Keep functions focused and small
- Write tests for new features

## Development Setup

```bash
git clone https://github.com/mrzhoong/go-api-scaffold.git
cd go-api-scaffold
go mod download
make run
```

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
