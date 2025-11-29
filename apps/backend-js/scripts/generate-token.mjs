import * as ed25519 from '@ucanto/principal/ed25519';
import * as Client from '@ucanto/client';

async function generateTestToken() {
  const alice = await ed25519.generate();
  const bob = await ed25519.generate();

  const delegation = await Client.delegate({
    issuer: alice,
    audience: bob,
    capabilities: [{
      can: 'store/add',
      with: 'storage:*',
    }],
    expiration: Math.floor(Date.now() / 1000) + 86400,
  });

  // Use archive() method which returns the CAR bytes
  const carBytes = await delegation.archive();
  console.log(Buffer.from(carBytes.ok).toString('base64'));
}

generateTestToken().catch(console.error);
