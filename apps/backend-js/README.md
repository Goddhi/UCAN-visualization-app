# UCAN Swagger Documentation Pattern

> **A reusable pattern for documenting UCAN-based APIs with OpenAPI/Swagger**

Add interactive API documentation to any UCAN-secured service with capability-based authorization, example token generation, and automatic validation.

---

A **complete reference implementation** showing how to integrate UCAN (User Controlled Authorization Networks) with OpenAPI/Swagger documentation.

This implementation provides:
- **Capability Registry** - Map endpoints to required UCAN capabilities
- **UCAN Validator Middleware** - Validate Bearer tokens with capability checking
- **Example Token Generator** - Auto-generate valid UCANs for testing
- **Interactive Swagger UI** - Try endpoints with UCAN tokens

---

### Prerequisites
```bash
node >= 18.0.0
yarn >= 1.22.0
```

### Installation
```bash
# Clone the repository
git clone https://github.com/Goddhi/UCAN-visualization-app
cd ucan-visualization-app/apps/backend-js

# Install dependencies
yarn install

# Start the server
yarn dev
```

### Visit the Documentation
```
http://localhost:8081/swagger
```

## How to Use This Pattern

### Step 1: Copy the Core Files

Copy these files to your project:
```bash
# Required files
src/middleware/ucanValidator.js
src/swagger/capabilities.js
src/swagger/config.js
src/swagger/examples.js
```

### Step 2: Define Your Service's Capabilities

**File:** `src/swagger/capabilities.js`
```javascript
export const capabilityRegistry = {
  'POST /your-endpoint': {
    capabilities: [{ can: 'your/action', with: 'your:resource' }],
    description: 'Your endpoint description'
  }
};
```

**Example for a file storage API:**
```javascript
export const capabilityRegistry = {
  'POST /files/upload': {
    capabilities: [{ can: 'store/add', with: 'storage:*' }],
    description: 'Upload files to storage'
  },
  'GET /files/:cid': {
    capabilities: [{ can: 'store/get', with: 'storage:*' }],
    description: 'Download files from storage'
  },
  'DELETE /files/:cid': {
    capabilities: [{ can: 'store/remove', with: 'storage:*' }],
    description: 'Delete files from storage'
  }
};
```

### Step 3: Generate Example Tokens

**File:** `src/swagger/examples.js`
```javascript
import * as ed25519 from '@ucanto/principal/ed25519';
import * as Client from '@ucanto/client';

export async function generateExampleTokens() {
  const issuer = await ed25519.generate();
  const audience = await ed25519.generate();

  // Generate tokens for YOUR capabilities
  const uploadToken = await Client.delegate({
    issuer,
    audience,
    capabilities: [{ can: 'store/add', with: 'storage:*' }],
    expiration: Math.floor(Date.now() / 1000) + 86400
  });
  
  const archive = await uploadToken.archive();
  return {
    upload: Buffer.from(archive.ok).toString('base64')
  };
}
```

### Step 4: Configure Swagger

**File:** `src/swagger/config.js`
```javascript
export const swaggerSpec = {
  openapi: '3.0.0',
  info: {
    title: 'Your API Name',
    version: '1.0.0',
    description: 'Your API description with UCAN info'
  },
  paths: {
    '/your-endpoint': {
      post: {
        summary: 'Your endpoint',
        description: '**Required capability:** `your/action` on `your:resource`',
        security: [{ UCANBearer: [] }],
        // ... rest of spec
      }
    }
  }
};
```

### Step 5: Add UCAN Validation to Your Server

**File:** `src/server.js`
```javascript
import express from 'express';
import swaggerUi from 'swagger-ui-express';
import { createUcanValidator } from './middleware/ucanValidator.js';
import { swaggerSpec } from './swagger/config.js';
import { generateExampleTokens } from './swagger/examples.js';

const app = express();

// Generate example tokens on startup
let exampleTokens = {};
generateExampleTokens().then(tokens => {
  exampleTokens = tokens;
});

// Enable UCAN validation
app.use(createUcanValidator());

// Swagger UI
app.use('/swagger', swaggerUi.serve, swaggerUi.setup(swaggerSpec));

// Example tokens endpoint
app.get('/api/examples/tokens', (req, res) => {
  res.json({ tokens: exampleTokens });
});

// Your routes here
app.use('/your-routes', yourRouter);

app.listen(3000);
```

---

## UCAN Validation Flow

### How It Works
```
1. Request arrives
   ↓
2. UCAN Validator Middleware extracts Bearer token
   ↓
3. Decodes UCAN using @ucanto/core/delegation
   ↓
4. Looks up required capabilities from registry
   ↓
5. Validates token has required capabilities
   ↓
6. Attaches delegation to req.ucan
   ↓
7. Proceeds to route handler OR returns 403
```

##  Testing

### Generate Test Tokens
```bash
# Start the server
yarn dev

### Test in Swagger UI

1. Visit `http://localhost:8081/swagger`
2. Click **"Authorize"** button
3. Enter: `Bearer <token-from-examples-endpoint>`
4. Click **"Authorize"** then **"Close"**
5. Try any endpoint with **"Try it out"**
