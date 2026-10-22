# Sourcey API reference configuration

This directory contains a [Sourcey](https://github.com/sourcey/sourcey)
configuration that builds a static, searchable API reference for Testify's five
public packages — `assert`, `require`, `mock`, `suite` and `http` — directly
from this repository's source.

Nothing is generated ahead of time: the reference is built from the tree, so it
can be regenerated for any commit and always matches the code.

## Building it

```bash
npm install sourcey@3.8.3

# Option 1: with the Go toolchain (live extraction)
npx sourcey godoc -m ../.. -p "./assert,./require,./mock,./suite,./http" -o godoc.json

# Option 2: from a committed snapshot (no Go toolchain needed)
# npx sourcey build -c docs/sourcey/sourcey.config.ts -o dist
```

`dist/` then contains a static site: a per-package API reference, a
machine-readable index (`llms.txt`), a sitemap, and Markdown twins of every
page.

## What it covers

The five public packages export 711 public entry points across 18 non-test
source files: 318 package-level functions, 355 methods, 39 types, plus the
package variables and constants. The generated reference renders all of them,
with each entry linked to the exact upstream line it came from.

## Why this is here

Testify's own documentation is a single `README.md` plus one package comment
per package. There is no per-symbol reference in this repository, and the
reference that most people read (`pkg.go.dev`) is auto-generated and reports no
coverage. A generated reference makes two things measurable:

- **Doc-comment coverage** — 710 of 714 exported symbols carry a doc comment
  (99.4%), so the prose is essentially complete.
- **Example coverage** — the whole public API ships 9 runnable `Example…`
  functions, all of them attached to the five `…AssertionFunc` types in
  `assert` and `require`. `mock`, `suite` and `http` have none.

That gap is the interesting one: a reference built from source shows exactly
which entry points a new user cannot find an example for.

## A published instance

An instance of this configuration, built from the pinned commit
`87a7b9d57689f6579db2da795ab8deeab29cb724` and published by the author of this
pull request, is live at <https://testify.anonymux.org/>. It is offered as a
reference implementation and as an opportunity to review the configuration
itself — if you would rather host it on a domain you control, the same
configuration and snapshot will produce the same site.

## Licence

The configuration in this directory is MIT licensed, like the rest of the
repository. The generated output contains Testify's own doc comments and
signatures and inherits the same licence.
