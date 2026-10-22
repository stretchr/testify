import { defineConfig, godoc } from "sourcey";

// Sourcey configuration for the Testify API reference.
//
// Sourcey (https://github.com/sourcey/sourcey) renders a static, searchable
// API reference straight from this repository's Go source. Nothing here is
// generated ahead of time: the site is built from the tree, so it can be
// regenerated for any commit.
//
//   npm install sourcey@3.8.3
//   npx sourcey godoc -m ../.. -p "./assert,./require,./mock,./suite,./http" -o godoc.json
//   npx sourcey build -c docs/sourcey/sourcey.config.ts -o dist
//
// The output in dist/ is plain static files and can be published anywhere that
// serves them.
export default defineConfig({
  name: "Testify API Reference",
  siteUrl: "https://testify.anonymux.org/",
  repo: "https://github.com/stretchr/testify",
  editBranch: "master",
  navigation: {
    tabs: [
      {
        tab: "API Reference",
        slug: "api",
        source: godoc({
          // The module root is the repository root, two levels above this file.
          module: "../..",
          packages: ["./assert", "./require", "./mock", "./suite", "./http"],
          // "auto" prefers the live Go toolchain and falls back to a committed
          // godoc.json snapshot when Go is not available.
          mode: "auto",
          includeTests: true,
        }),
      },
    ],
  },
});
