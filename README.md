# UCAN Visualization Web App

UCAN Visualization Web App is a comprehensive developer tool for parsing, validating, and visualizing UCAN (User Controlled Authorization Networks) delegation chains. This application addresses the core challenge of debugging decentralized authorization by providing visual feedback and detailed error reporting for UCAN tokens.

## Quick Start

### Prerequisites

- **Node.js** >= 18
- **Yarn** 1.22.22 (or npm/pnpm)
- **Go** >= 1.21 (for backend)
- **Air** (Go hot reload tool)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd UCAN
   ```

2. **Install dependencies**
   ```bash
   yarn install
   ```

3. **Install Go tools** (required for backend hot reload)
   ```bash
   make install-tools
   ```
   This installs Air for Go hot reloading.

### Running the Application

**Development Mode (Both apps)**
```bash
yarn dev
```
Runs both frontend (port 3000) and backend in parallel.

```bash
yarn dev:reload
```

Runs both frontend (port 3000) and backend (using Air) in parallel, with hot reload on the backend.

**Frontend Only**
```bash
yarn dev-frontend
```
Starts Next.js frontend on http://localhost:3000

**Backend Only**
```bash
yarn dev-backend
```
Starts Go backend with Air hot reloading.

## Available Commands

### Root Commands (from project root)

| Command | Description |
|---------|-------------|
| `yarn dev` | Run both frontend and backend in parallel |
| `yarn dev-frontend` | Run only the frontend application |
| `yarn dev-backend` | Run only the backend with Air hot reload |
| `yarn dev:reload` | Run both apps with hot reload enabled |
| `yarn build` | Build all applications |
| `yarn lint` | Lint all applications |
| `yarn format` | Format code with Prettier |
| `yarn check-types` | Type-check all TypeScript files |

### Backend Commands (from apps/backend)

| Command | Description |
|---------|-------------|
| `yarn dev` | Run backend server (direct Go run) |
| `yarn dev:reload` | Run backend with Air hot reload |
| `yarn build` | Build backend binary to bin/server |
| `yarn lint` | Run golangci-lint |
| `yarn test` | Run Go tests |

### Frontend Commands (from apps/frontend)

| Command | Description |
|---------|-------------|
| `yarn dev` | Run frontend with Turbopack (fast refresh) |
| `yarn dev:reload` | Run frontend in standard mode |
| `yarn build` | Build production-ready frontend |
| `yarn start` | Start production server |
| `yarn lint` | Run ESLint |
| `yarn check-types` | Type-check TypeScript |

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make install-tools` | Install Air for Go hot reloading |

## Project Structure

```
UCAN/
├── apps/
│   ├── backend/          # Go backend API
│   │   ├── cmd/          # Application entry points
│   │   ├── internal/     # Internal packages
│   │   │   ├── api/      # HTTP handlers and routing
│   │   │   ├── config/   # Configuration management
│   │   │   ├── models/   # Data models
│   │   │   └── services/ # Business logic services
│   │   ├── pkg/          # Public packages
│   │   └── test/         # Tests and fixtures
│   └── frontend/         # Next.js frontend
│       └── app/          # Next.js 15 app directory
├── packages/
│   ├── eslint-config/    # Shared ESLint configurations
│   ├── typescript-config/# Shared TypeScript configurations
│   └── ui/               # Shared UI components
└── turbo.json           # Turborepo configuration
```

### Features

**Delegation Chain Visualizer**

Visual tree/graph showing who delegated what to whom
Interactive exploration of trust relationships
Color-coded validity indicators (valid/expired/invalid)
See the complete authorization path from root to end user

**Invocation Inspector**

Paste a UCAN invocation and see exactly what it's trying to do
View the capability being invoked (resource, ability, caveats)
Inspect all attached proofs in the invocation
Understand invocation metadata (issuer, audience, timestamps)

**Capability Breakdown**

Parse and display resources, abilities, and attenuations in human-readable format
See technical capability structure translated to plain English
Understand what permissions are granted and what restrictions apply
Identify how capabilities were narrowed through delegation (attenuation)

**Proof Chain Validator**

Check if a UCAN chain is valid and where it breaks if invalid
Cryptographic signature verification
Time-based validation (expiration, not-before)
Capability escalation detection
Clear root cause analysis with fix suggestions
