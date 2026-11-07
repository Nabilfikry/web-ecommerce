# E-Commerce Monorepo

Monorepo full-stack e-commerce architecture powered by **Turborepo** + **Bun** + **Go**.

## Structure

```tree
e-commerce/
├── apps/
│   ├── web/                # Frontend: Svelte 5, SvelteKit, Tailwind CSS v4, Storybook
│   └── backend-service/    # Backend: Go HTTP service (cmd, internal, pkg, api)
├── packages/
│   └── shared-types/       # Shared TypeScript types & interfaces
├── turbo.json              # Turborepo task pipeline configuration
├── go.work                 # Go workspace multi-module root
├── package.json            # Root workspace scripts & dependencies
└── bun.lock                # Bun lockfile
```

## Quick Start

### Prerequisites
- [Bun](https://bun.sh) (v1.2+)
- [Go](https://go.dev) (v1.22+)

### Installation

```bash
bun install
```

### Development

Run all applications (web frontend + Go backend) concurrently with caching:

```bash
bun run dev
```

Or run services individually:

```bash
# Frontend only (http://localhost:5173)
bun run dev:web

# Backend only (http://localhost:8080)
bun run dev:backend
```

### Build

```bash
bun run build
```

### Lint & Quality Checks

```bash
bun run lint
bun run check
bun run test
```

