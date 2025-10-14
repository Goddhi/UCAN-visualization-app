# UCAN Visualizer API Documentation

## Overview


The UCAN Visualizer API provides endpoints for parsing, validating, and visualizing UCAN (User Controlled Authorization Network) delegation tokens.

**Supported Input Formats:**
- Base64-encoded CAR (Content Addressed aRchive) - via JSON
- Binary CAR file - via multipart file upload
- Hex-encoded string - via JSON with format hint

### Start server
```go run ./cmd/server/```

### Testing
Get a Test Token
Generate a test UCAN token:
```go run scripts/generate-token.go```
This will output a base64-encoded token you can use for testing.

## Endpoints

### 1. Parse Delegation

Parse a UCAN token from JSON and extract delegation information.

**Endpoint:** `POST /api/parse/delegation`

**Request Headers:**

Content-Type: application/json

**Request Body:**
```json
{
  "token": "Y0c5WkM3RD...",
  "format": "base64"  // optional: "base64", "hex", or "auto" (default)
}
```
**Success Response: 200 OK**

```
json{
  "issuer": "did:key:z6MkfXPzGKT1YKMZsW4797KMfgUnx2TTwkQthyx25jetV93b",
  "audience": "did:key:z6Mkt3Q73FauA4TWE7RBUUURJx239N55s6wyB9F15zWVSvoV",
  "capabilities": [
    {
      "with": "storage:alice/*",
      "can": "store/add",
      "nb": {}
    }
  ],
  "proofs": [],
  "expiration": "2025-10-14T22:37:18Z",
  "notBefore": "2025-10-13T22:37:18Z",
  "signature": {
    "algorithm": "EdDSA"
  },
  "cid": "bafyreihs233ufeujd5pjz2res5v2ymtmfypfvg7jprvcefdg7d3akyayty"
}
```

**Error Responses:**
400 Bad Request - Invalid token format or missing token
422 Unprocessable Entity - Valid format but failed to parse UCAN

#### Parse Delegation (File)
Parse a UCAN token from an uploaded CAR file.
Endpoint: POST /api/parse/delegation/file
Request Headers:
Content-Type: multipart/form-data
Request Body:

Field name: file
File types: .car, .ucan, .cbor, or no extension
Max size: 10 MB

Success Response: 200 OK
Same as /api/parse/delegation

**Error Responses:**
400 Bad Request - No file provided or invalid file type
422 Unprocessable Entity - Failed to parse delegation

### Validate Chain
Validate a UCAN delegation token, checking time bounds and structural integrity.
Endpoint: POST /api/validate/chain
Request Headers:
Content-Type: application/json
Request Body:
```
json{
  "token": "Y0c5WkM3RD...",
  "format": "base64"  // optional
}
Success Response: 200 OK
json{
  "valid": true,
  "chain": [
    {
      "level": 0,
      "cid": "bafyreihs233ufeujd5pjz2res5v2ymtmfypfvg7jprvcefdg7d3akyayty",
      "issuer": "did:key:z6MkfXPzGKT1YKMZsW4797KMfgUnx2TTwkQthyx25jetV93b",
      "audience": "did:key:z6Mkt3Q73FauA4TWE7RBUUURJx239N55s6wyB9F15zWVSvoV",
      "capability": {
        "with": "storage:alice/*",
        "can": "store/add",
        "nb": {}
      },
      "expiration": "2025-10-14T22:37:18Z",
      "notBefore": "2025-10-13T22:37:18Z",
      "valid": true,
      "issues": [
        {
          "type": "expiring_soon",
          "message": "UCAN expires in 23h 40m",
          "severity": "warning",
          "context": {
            "expiration": "2025-10-14T22:37:18Z",
            "time_remaining": "23h40m0s"
          }
        }
      ]
    }
  ],
  "rootCause": null,
  "summary": {
    "totalLinks": 1,
    "validLinks": 1,
    "invalidLinks": 0,
    "warningCount": 1
  }
}
```

**Validation Checks:**
Expiration time (is the UCAN expired?)
Not-before time (is the UCAN active yet?)
Structural integrity (valid capabilities, proofs)

**Error Responses:**
400 Bad Request - Invalid token format
500 Internal Server Error - Validation failed

#### Validate Chain (File)
Validate a UCAN token from an uploaded CAR file.
Endpoint: POST /api/validate/chain/file
Request: Same as /api/parse/delegation/file
Success Response: Same as /api/validate/chain

### Generate Graph
Generate graph visualization data for a UCAN delegation.
Endpoint: POST /api/graph/delegation
Request Headers:
Content-Type: application/json
Request Body:
json{
  "token": "Y0c5WkM3RD...",
  "format": "base64"  // optional
}
Success Response: 200 OK
```
json{
  "nodes": [
    {
      "id": "did:key:z6MkfXPzGKT1YKMZsW4797KMfgUnx2TTwkQthyx25jetV93b",
      "label": "did:key:z6MkfX...93b",
      "type": "root",
      "metadata": {
        "fullDid": "did:key:z6MkfXPzGKT1YKMZsW4797KMfgUnx2TTwkQthyx25jetV93b",
        "role": "issuer"
      }
    },
    {
      "id": "did:key:z6Mkt3Q73FauA4TWE7RBUUURJx239N55s6wyB9F15zWVSvoV",
      "label": "did:key:z6Mkt3...voV",
      "type": "leaf",
      "metadata": {
        "delegationCID": "bafyreihs233ufeujd5pjz2res5v2ymtmfypfvg7jprvcefdg7d3akyayty",
        "expiration": 1760477838,
        "fullDid": "did:key:z6Mkt3Q73FauA4TWE7RBUUURJx239N55s6wyB9F15zWVSvoV",
        "notBefore": 1760391438,
        "role": "audience"
      }
    }
  ],
  "edges": [
    {
      "source": "did:key:z6MkfXPzGKT1YKMZsW4797KMfgUnx2TTwkQthyx25jetV93b",
      "target": "did:key:z6Mkt3Q73FauA4TWE7RBUUURJx239N55s6wyB9F15zWVSvoV",
      "capability": {
        "with": "storage:alice/*",
        "can": "store/add",
        "nb": {
          "_hasCaveats": true
        }
      },
      "valid": true,
      "label": "store/add on storage:alice/*"
    }
  ]
}
```
**Graph Structure:**
**Nodes**: Represent principals (DIDs)

**root**: Original issuer
**leaf**: Final audience
**intermediate**: Middle of chain (for multi-level delegations)


**Edges**: Represent delegations with capabilities