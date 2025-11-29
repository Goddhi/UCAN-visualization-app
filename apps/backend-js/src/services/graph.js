export class GraphService {
  constructor(parserService) {
    this.parser = parserService;
  }

  async generateDelegationGraph(tokenBytes) {
    const chain = await this.parser.parseDelegationChain(tokenBytes);
    const { nodes, edges } = this.#buildGraph(chain);
    const chainInfo = this.#buildChainInfo(chain);

    return { nodes, edges, chain: chainInfo };
  }

  async generateInvocationGraph(tokenBytes) {
    const invocation = await this.parser.parseInvocation(tokenBytes);
    const chain = await this.parser.parseDelegationChain(tokenBytes);
    const { nodes, edges } = this.#buildGraph(chain);
    const chainInfo = this.#buildChainInfo(chain);

    return {
      nodes,
      edges,
      chain: chainInfo,
      invocation,
      isInvocation: invocation.isInvocation
    };
  }

  #buildGraph(chain) {
    const nodesMap = new Map();
    const edges = [];
    const maxLevel = Math.max(...chain.map(d => d.level));

    for (const del of chain) {
      if (!nodesMap.has(del.issuer)) {
        nodesMap.set(del.issuer, {
          id: del.issuer,
          label: this.#shortenDID(del.issuer),
          type: del.level === 0 ? 'root' : 'intermediate',
          level: del.level,
          metadata: { fullDID: del.issuer, role: 'delegator' }
        });
      }

      if (!nodesMap.has(del.audience)) {
        nodesMap.set(del.audience, {
          id: del.audience,
          label: this.#shortenDID(del.audience),
          type: del.level === maxLevel ? 'leaf' : 'intermediate',
          level: del.level,
          metadata: { fullDID: del.audience, role: 'delegatee' }
        });
      }

      for (const cap of del.capabilities) {
        edges.push({
          source: del.issuer,
          target: del.audience,
          capability: cap,
          label: `${cap.can} on ${cap.with}`,
          valid: true,
          level: del.level,
          type: 'delegation',
          metadata: { cid: del.cid }
        });
      }
    }

    return {
      nodes: Array.from(nodesMap.values()),
      edges
    };
  }

  #buildChainInfo(chain) {
    if (chain.length === 0) return {};

    const maxLevel = Math.max(...chain.map(d => d.level));
    const leafCIDs = chain.filter(d => d.level === maxLevel).map(d => d.cid);

    return {
      totalLevels: maxLevel + 1,
      isComplete: true,
      rootCID: chain[0].cid,
      leafCIDs,
      principals: [],
      timeline: []
    };
  }

  #shortenDID(did) {
    if (did.length <= 20) return did;
    const parts = did.split(':');
    if (parts.length < 3) return did.slice(0, 20) + '...';
    
    const identifier = parts[2];
    if (identifier.length > 16) {
      return `did:${parts[1]}:${identifier.slice(0, 8)}...${identifier.slice(-8)}`;
    }
    return did;
  }
}
