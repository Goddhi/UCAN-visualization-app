import * as Delegation from '@ucanto/core/delegation';
import * as Block from 'multiformats/block';
import * as dagCBOR from '@ipld/dag-cbor';
import { sha256 } from 'multiformats/hashes/sha2';

export class ParserService {
  async parseDelegation(tokenBytes) {
    try {
      const result = await Delegation.extract(tokenBytes);
      const delegation = result.ok || result;
      return this.delegationToResponse(delegation, 0);
    } catch (error) {
      throw new Error(`Failed to parse delegation: ${error.message}`);
    }
  }

  async parseDelegationChain(tokenBytes) {
    try {
      const result = await Delegation.extract(tokenBytes);
      const delegation = result.ok || result;
      const chain = [];
      
      const root = this.delegationToResponse(delegation, 0);
      chain.push(root);

      if (delegation.proofs && delegation.proofs.length > 0) {
        const proofDelegations = await this.parseProofs(
          delegation.proofs, 
          delegation.blocks,
          1
        );
        chain.push(...proofDelegations);
      }

      return chain;
    } catch (error) {
      throw new Error(`Failed to parse delegation chain: ${error.message}`);
    }
  }

  async parseInvocation(tokenBytes) {
    const delegation = await this.parseDelegation(tokenBytes);
    const invocationAnalysis = this.analyzeInvocation(delegation);
    const capabilityAnalysis = this.analyzeCapabilities(delegation.capabilities);

    let task = null;
    if (invocationAnalysis.isInvocation) {
      task = {
        action: invocationAnalysis.primaryAction,
        resource: invocationAnalysis.targetResource,
        constraints: invocationAnalysis.constraints,
        issuer: delegation.issuer,
        target: delegation.audience,
        taskType: invocationAnalysis.taskType,
        permissions: invocationAnalysis.requiredPermissions
      };
    }

    return {
      delegation,
      isInvocation: invocationAnalysis.isInvocation,
      task,
      invocationAnalysis,
      capabilityAnalysis
    };
  }

  delegationToResponse(delegation, level) {
    const caps = delegation.data?.capabilities || delegation.capabilities || [];
    
    const capabilities = caps.map(cap => ({
      with: cap.with,
      can: cap.can,
      nb: cap.nb || {},
      category: this.categorizeCapability(cap.can)
    }));

    const proofs = (delegation.proofs || []).map((proof, index) => ({
      cid: proof.toString ? proof.toString() : String(proof),
      index,
      type: 'delegation'
    }));

    const exp = delegation.data?.expiration || delegation.expiration;
    const nbf = delegation.data?.notBefore || delegation.notBefore;
    
    const expiration = exp ? new Date(exp * 1000).toISOString() : null;
    const notBefore = nbf ? new Date(nbf * 1000).toISOString() : null;

    const issuer = delegation.issuer?.did ? delegation.issuer.did() : String(delegation.issuer || 'unknown');
    const audience = delegation.audience?.did ? delegation.audience.did() : String(delegation.audience || 'unknown');
    const cid = delegation.cid?.toString ? delegation.cid.toString() : String(delegation.cid || 'unknown');

    return {
      issuer,
      audience,
      capabilities,
      proofs,
      expiration,
      notBefore,
      facts: delegation.data?.facts || delegation.facts || [],
      nonce: delegation.data?.nonce || delegation.nonce || null,
      signature: { algorithm: 'EdDSA' },
      cid,
      level
    };
  }

  async parseProofs(proofLinks, blocks, level) {
    const proofs = [];

    for (const link of proofLinks) {
      try {
        const blockBytes = blocks.get(link);
        if (!blockBytes) continue;

        const block = await Block.decode({
          bytes: blockBytes,
          codec: dagCBOR,
          hasher: sha256
        });

        const proofDelegation = await this.blockToDelegation(block, blocks);
        const parsed = this.delegationToResponse(proofDelegation, level);
        proofs.push(parsed);

        if (proofDelegation.proofs && proofDelegation.proofs.length > 0) {
          const nested = await this.parseProofs(
            proofDelegation.proofs,
            blocks,
            level + 1
          );
          proofs.push(...nested);
        }
      } catch (error) {
        console.warn(`Failed to parse proof ${link}:`, error.message);
      }
    }

    return proofs;
  }

  async blockToDelegation(block, blocks) {
    const value = block.value;
    return {
      cid: block.cid,
      issuer: { did: () => value.iss },
      audience: { did: () => value.aud },
      data: {
        capabilities: value.att || [],
        expiration: value.exp,
        notBefore: value.nbf,
        facts: value.fct || [],
        nonce: value.nnc,
      },
      proofs: value.prf || [],
      blocks
    };
  }

  analyzeInvocation(delegation) {
    const analysis = {
      isInvocation: false,
      hasInvokeCapability: false,
      taskType: 'delegation',
      primaryAction: '',
      targetResource: '',
      invokePatterns: [],
      requiredPermissions: [],
      constraints: {}
    };

    if (delegation.issuer !== delegation.audience) {
      analysis.isInvocation = true;
    }

    for (const cap of delegation.capabilities) {
      if (this.isInvokeCapability(cap.can)) {
        analysis.hasInvokeCapability = true;
        analysis.taskType = 'invocation';
        analysis.invokePatterns.push(cap.can);
        analysis.primaryAction = cap.can;
        analysis.targetResource = cap.with;
      }

      if (cap.can) {
        analysis.requiredPermissions.push(cap.can);
      }

      if (cap.nb) {
        Object.assign(analysis.constraints, cap.nb);
      }
    }

    if (!analysis.primaryAction && delegation.capabilities.length > 0) {
      analysis.primaryAction = delegation.capabilities[0].can;
      analysis.targetResource = delegation.capabilities[0].with;
    }

    return analysis;
  }

  analyzeCapabilities(capabilities) {
    const analysis = {
      categories: {},
      totalCount: capabilities.length,
      invokeCount: 0,
      delegateCount: 0,
      permissions: [],
      resources: []
    };

    for (const cap of capabilities) {
      const category = cap.category || 'unknown';
      
      if (!analysis.categories[category]) {
        analysis.categories[category] = [];
      }
      analysis.categories[category].push(cap);

      if (this.isInvokeCapability(cap.can)) {
        analysis.invokeCount++;
      } else {
        analysis.delegateCount++;
      }

      if (cap.can) analysis.permissions.push(cap.can);
      if (cap.with) analysis.resources.push(cap.with);
    }

    return analysis;
  }

  isInvokeCapability(capability) {
    const invokePatterns = ['invoke', 'execute', 'run', 'call', 'perform'];
    return invokePatterns.some(pattern => 
      capability.toLowerCase().includes(pattern)
    );
  }

  categorizeCapability(capability) {
    if (capability.startsWith('store')) return 'storage';
    if (capability.startsWith('space')) return 'space';
    if (capability.startsWith('upload')) return 'upload';
    if (this.isInvokeCapability(capability)) return 'invocation';
    if (capability.startsWith('blob')) return 'blob';
    if (capability.startsWith('index')) return 'index';
    return 'general';
  }
}
