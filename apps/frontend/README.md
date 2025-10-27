# UCAN Visualizer Frontend

Interactive web interface for parsing, validating, and visualizing UCAN delegation chains.

## Prerequisites

- Node.js 18+ and npm/yarn
- Backend API running (see `apps/backend`)

## Setup

1. Install dependencies:
```bash
yarn install
```

2. Configure environment variables:
```bash
cp .env.example .env.local
```

Edit `.env.local` and set the backend API URL:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Development

Start the development server:
```bash
yarn dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Features

- **UCAN Token Parser**: Parse UCAN tokens from text input or file upload (.car, .ucan, .cbor)
- **Interactive Graph Visualizer**: Visualize delegation chains as interactive flowcharts
- **Delegation Details**: View issuer, audience, capabilities, and expiration info
- **Real-time Validation**: Validate UCAN chains with the backend API

## API Integration

The frontend connects to the backend API for:
- `/api/parse/delegation` - Parse UCAN tokens
- `/api/parse/delegation/file` - Parse UCAN from uploaded files
- `/api/validate/chain` - Validate delegation chains
- `/api/graph/delegation` - Generate graph data

## Project Structure

```
app/
├── components/       # React components
├── graph/           # Graph visualization page
├── settings/        # Settings page
├── lib/
│   ├── api/        # API client and types
│   └── utils/      # Utility functions
└── page.tsx        # Home page
```

## Build

Build for production:
```bash
yarn build
```

Start production server:
```bash
yarn start
```
