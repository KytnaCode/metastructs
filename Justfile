export CGO_ENABLED := '0'

check:
  govulncheck ./...

lint:
  golangci-lint run ./...

format:
  golangci-lint fmt ./...

precommit: lint format check
