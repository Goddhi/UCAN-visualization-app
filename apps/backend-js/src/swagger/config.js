/**
 * Simple Swagger configuration with UCAN security
 */
export const swaggerSpec = {
    openapi: '3.0.0',
    info: {
      title: 'UCAN Visualization API',
      version: '1.0.0',
      description: `
  # UCAN-Secured API
  
  This API validates UCAN tokens for authorization.
  
  ## How to Use
  
  1. Get a UCAN token from /api/examples/tokens
  2. Click "Authorize" button below
  3. Enter: \`Bearer <your-token>\`
  4. Try the endpoints
  
  ## Example Tokens
  
  Visit http://localhost:8081/api/examples/tokens to get fresh tokens.
  
  Each endpoint requires specific capabilities:
  - **Parse endpoints**: \`ucan/parse\` on \`api:*\`
  - **Validate endpoints**: \`ucan/validate\` on \`api:*\`
  - **Graph endpoints**: \`ucan/graph\` on \`api:*\`
  - **Admin token**: All capabilities
      `
    },
    servers: [
      { url: 'http://localhost:8081', description: 'Development' }
    ],
    components: {
      securitySchemes: {
        UCANBearer: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'UCAN',
          description: 'UCAN token with required capabilities. Get tokens from /api/examples/tokens'
        }
      },
      schemas: {
        TokenInput: {
          type: 'object',
          required: ['token'],
          properties: {
            token: {
              type: 'string',
              format: 'byte',
              description: 'Base64-encoded UCAN token in CAR format'
            }
          }
        },
        Error: {
          type: 'object',
          properties: {
            error: { type: 'string' },
            message: { type: 'string' }
          }
        }
      }
    },
    security: [{ UCANBearer: [] }],
    paths: {
      '/health': {
        get: {
          tags: ['System'],
          summary: 'Health check',
          description: 'Check if the service is running',
          security: [],
          responses: {
            '200': { 
              description: 'Service is healthy',
              content: {
                'application/json': {
                  schema: {
                    type: 'object',
                    properties: {
                      status: { type: 'string' },
                      service: { type: 'string' },
                      timestamp: { type: 'string' },
                      version: { type: 'string' }
                    }
                  }
                }
              }
            }
          }
        }
      },
      '/api/examples/tokens': {
        get: {
          tags: ['System'],
          summary: 'Get example UCAN tokens',
          description: 'Returns example tokens for testing each endpoint',
          security: [],
          responses: {
            '200': {
              description: 'Example tokens',
              content: {
                'application/json': {
                  schema: {
                    type: 'object',
                    properties: {
                      message: { type: 'string' },
                      tokens: { type: 'object' },
                      usage: { type: 'object' }
                    }
                  }
                }
              }
            }
          }
        }
      },
      '/api/parse/delegation': {
        post: {
          tags: ['Parse'],
          summary: 'Parse UCAN delegation',
          description: '**Required capability:** `ucan/parse` on `api:*`\n\nExtracts delegation information from a UCAN token.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Successfully parsed delegation' },
            '401': { 
              description: 'Missing UCAN token',
              content: {
                'application/json': {
                  schema: { $ref: '#/components/schemas/Error' }
                }
              }
            },
            '403': { 
              description: 'Insufficient capabilities',
              content: {
                'application/json': {
                  schema: { $ref: '#/components/schemas/Error' }
                }
              }
            },
            '422': {
              description: 'Failed to parse delegation',
              content: {
                'application/json': {
                  schema: { $ref: '#/components/schemas/Error' }
                }
              }
            }
          }
        }
      },
      '/api/parse/chain': {
        post: {
          tags: ['Parse'],
          summary: 'Parse delegation chain',
          description: '**Required capability:** `ucan/parse` on `api:*`\n\nParses a full UCAN delegation chain including all proofs.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Successfully parsed chain' },
            '403': { description: 'Insufficient capabilities' },
            '422': { description: 'Failed to parse chain' }
          }
        }
      },
      '/api/parse/invocation': {
        post: {
          tags: ['Parse'],
          summary: 'Parse UCAN invocation',
          description: '**Required capability:** `ucan/parse` on `api:*`\n\nAnalyzes a UCAN invocation and extracts task information.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Successfully analyzed invocation' },
            '403': { description: 'Insufficient capabilities' },
            '422': { description: 'Failed to parse invocation' }
          }
        }
      },
      '/api/validate/chain': {
        post: {
          tags: ['Validate'],
          summary: 'Validate delegation chain',
          description: '**Required capability:** `ucan/validate` on `api:*`\n\nValidates a UCAN delegation chain and identifies any issues.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Validation result' },
            '403': { description: 'Insufficient capabilities' },
            '500': { description: 'Validation failed' }
          }
        }
      },
      '/api/graph/delegation': {
        post: {
          tags: ['Graph'],
          summary: 'Generate delegation graph',
          description: '**Required capability:** `ucan/graph` on `api:*`\n\nGenerates a node-edge graph for visualizing the delegation.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Graph data' },
            '403': { description: 'Insufficient capabilities' },
            '422': { description: 'Failed to generate graph' }
          }
        }
      },
      '/api/graph/invocation': {
        post: {
          tags: ['Graph'],
          summary: 'Generate invocation graph',
          description: '**Required capability:** `ucan/graph` on `api:*`\n\nGenerates an enhanced graph with invocation markers.',
          security: [{ UCANBearer: [] }],
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: { $ref: '#/components/schemas/TokenInput' }
              }
            }
          },
          responses: {
            '200': { description: 'Invocation graph data' },
            '403': { description: 'Insufficient capabilities' },
            '422': { description: 'Failed to generate graph' }
          }
        }
      }
    }
  };