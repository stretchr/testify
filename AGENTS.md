# AGENTS.md

Context file for AI agents working on testify.

**Dual Format**: This file combines Category A (Operations Manual) and Category B (Context Guide) for comprehensive agent guidance.

## Project Overview

testify is a Go project using Unknown.

**Key Info:**
- **Primary Language:** Go
- **Build System:** Unknown
- **Test Framework:** Go testing
- **Total Files:** 91
- **Test Files:** 24
- **AI Readiness Score:** 82/100 (AI-Native-Plus)

---

## 🚨 AI Policy & Operations

Extracted from CONTRIBUTING.md - operational constraints and procedures.

### AI Policy

- the community. A few general guidelines are listed here:

### Key Requirements

- need to know to get up to speed on the project.
- need to enforce certain practices and patterns to keep the package easy for others to understand.

### Development Procedures

- So you'd like to contribute to Testify? First of all, thank you! Testify is widely used, so each
- The Testify maintainers generally attempt to follow widely accepted practices within the Golang
- that we should try hard to use meaningful, descriptive function names, avoid breaking the API
- Quality isn't an afterthought* - Testify is a testing library, so it seems reasonable that we
- should have a decent test suite. This is doubly important because a bug in Testify doesn't just mean



## 🏗️ Architecture & Context Guide

This section provides architectural context and agent-understanding for the codebase.

### Prerequisites

- **Go:** 1.18+ (or applicable language version)
- **Package Manager:** go modules
- **Test Runner:** Go testing

### Environment Requirements

- **Go:** 1.17+ (from `go.mod`)
  - GCC required for CGo/SQLite compilation
- **Package Manager:** go modules


### Project Structure

```
testify/
├── src/                  # Source code
├── tests/                # Test suite (24 files)
└── README.md             # Project documentation
```

### Architecture Overview

#### Key Components
- **Main Entry:** main.go, main.go
- **Test Suite:** 24 test files
- **Build Configuration:** Standard

#### Design Principles

1. **Modularity** - Code organized by functionality with clear separation of concerns
2. **Testability** - Comprehensive test coverage across critical paths
3. **Clarity** - Explicit naming and structure for AI agent understanding
4. **Consistency** - Uniform patterns and conventions throughout codebase
5. **Maintainability** - Well-documented code with clear intent

### Directory Map

| Directory | Purpose |
|-----------|----------|
| `src/` or project root | Main source code |
| `tests/` or `test/` | Test suite |


### Development Workflow

#### Initial Setup

```bash
git clone https://github.com/YOUR_ORG/testify.git
cd testify
go mod download
```

#### Development Commands

**Running Tests:**
```bash
go build ./...            # Build project
go test ./...             # Run all tests
go test -v ./...          # Verbose test output
golangci-lint run         # Lint (if installed)
```

#### Code Quality
```bash
gofmt -w .                # Format code
go vet ./...              # Vet (static analysis)
```

### Code Style & Conventions

- **Naming:** Use Go conventions (snake_case for functions, PascalCase for classes)
- **Type Hints:** Yes (strongly encouraged)
- **Error Handling:** Yes - handle errors at boundaries; let exceptions propagate when another layer owns recovery
- **Logging:** Yes
- **Testing:** Yes - write tests alongside code changes

### Testing Strategy

**Framework:** Go testing
**Test Files:** 24 found

Before committing:
1. Run the full test suite: `go test ./...`
2. Ensure all tests pass: `go test -v ./...`
3. Run linter: `golangci-lint run`
4. Format code: `gofmt -w .`

### Writing Documentation

When updating docs:
1. Always include explanatory text before code snippets
2. Describe *why* and *what* before showing *how*
3. Keep sections focused on a single concept
4. Use clear, concrete examples

## Known Gotchas & Warnings

- provides meaningful benefits. For example, don't add functions because they might be useful to
- that we should try hard to use meaningful, descriptive function names, avoid breaking the API

### Contributing Guidelines

This project has a detailed contribution guide at **`CONTRIBUTING.md`**.

**Key Requirements:**
- Review the contribution guide for all requirements
- Follow established patterns in the codebase
- Ensure alignment with project's contribution policies

### Common Patterns

When contributing to this project:
1. Read existing code in the area you're modifying
2. Follow the established patterns and style
3. Write tests for new functionality
4. Use clear, descriptive variable and function names
5. Add docstrings for public APIs
6. Update tests when changing behavior

### What We Value

✅ Well-tested code with clear intent
✅ Consistent code style and naming conventions
✅ Code that is easy for AI agents to understand
✅ Clear, descriptive commit messages
✅ Modular, reusable components
✅ Comprehensive documentation

### What We Avoid

❌ Large functions doing multiple things
❌ Commented-out dead code
❌ Inconsistent naming or patterns
❌ Unclear error messages
❌ Unexplained magic numbers or strings
❌ Skipped tests or test TODOs

### AI Readiness Dimensions (Scoring)

This project is evaluated across 8 dimensions:

1. **Architecture** (20/100) - Code organization and modularity
2. **Testing** (15/100) - Test coverage and quality
3. **Dependencies** (6/100) - Dependency management
4. **Conventions** (8/100) - Consistent patterns
5. **Entry Points** (10/100) - Clear main/start locations
6. **Security** (10/100) - Input validation and error handling
7. **Build** (5/100) - Clear build/setup instructions
8. **Documentation** (8/100) - Code and project documentation

### Next Steps

Before making changes:
1. Read relevant source files to understand the existing code
2. Look at existing tests for similar functionality
3. Follow the patterns you see in the codebase
4. Write tests for your changes
5. Run `pytest` to verify nothing breaks
6. Run code quality checks: `ruff check . && mypy .`
7. Format your code: `ruff format .`

---

*Generated by Braxis - keeping AI agents in sync with your code*
