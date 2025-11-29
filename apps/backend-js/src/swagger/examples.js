import * as ed25519 from '@ucanto/principal/ed25519';
import * as Client from '@ucanto/client';

/**
 * Generate example UCAN tokens for documentation
 */
export async function generateExampleTokens() {
  const issuer = await ed25519.generate();
  const audience = await ed25519.generate();

  const tokens = {};

  // Parse token
  const parseDelegation = await Client.delegate({
    issuer,
    audience,
    capabilities: [{ can: 'ucan/parse', with: 'api:*' }],
    expiration: Math.floor(Date.now() / 1000) + 86400
  });
  const parseArchive = await parseDelegation.archive();
  tokens.parse = Buffer.from(parseArchive.ok).toString('base64');

  // Validate token
  const validateDelegation = await Client.delegate({
    issuer,
    audience,
    capabilities: [{ can: 'ucan/validate', with: 'api:*' }],
    expiration: Math.floor(Date.now() / 1000) + 86400
  });
  const validateArchive = await validateDelegation.archive();
  tokens.validate = Buffer.from(validateArchive.ok).toString('base64');

  // Graph token
  const graphDelegation = await Client.delegate({
    issuer,
    audience,
    capabilities: [{ can: 'ucan/graph', with: 'api:*' }],
    expiration: Math.floor(Date.now() / 1000) + 86400
  });
  const graphArchive = await graphDelegation.archive();
  tokens.graph = Buffer.from(graphArchive.ok).toString('base64');

  // Admin token (all capabilities)
  const adminDelegation = await Client.delegate({
    issuer,
    audience,
    capabilities: [
      { can: 'ucan/parse', with: 'api:*' },
      { can: 'ucan/validate', with: 'api:*' },
      { can: 'ucan/graph', with: 'api:*' }
    ],
    expiration: Math.floor(Date.now() / 1000) + 86400
  });
  const adminArchive = await adminDelegation.archive();
  tokens.admin = Buffer.from(adminArchive.ok).toString('base64');

  return tokens;
}