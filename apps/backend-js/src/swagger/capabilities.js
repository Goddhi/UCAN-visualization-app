/**
 * Capability Registry - Maps endpoints to required UCAN capabilities
 */
export const capabilityRegistry = {
    // Parse endpoints - require 'parse' capability
    'POST /api/parse/delegation': {
      capabilities: [{ can: 'ucan/parse', with: 'api:*' }],
      description: 'Parse UCAN delegation tokens'
    },
    'POST /api/parse/chain': {
      capabilities: [{ can: 'ucan/parse', with: 'api:*' }],
      description: 'Parse UCAN delegation chains'
    },
    'POST /api/parse/invocation': {
      capabilities: [{ can: 'ucan/parse', with: 'api:*' }],
      description: 'Parse UCAN invocations'
    },
    
    // Validate endpoints - require 'validate' capability
    'POST /api/validate/chain': {
      capabilities: [{ can: 'ucan/validate', with: 'api:*' }],
      description: 'Validate UCAN delegation chains'
    },
    
    // Graph endpoints - require 'graph' capability
    'POST /api/graph/delegation': {
      capabilities: [{ can: 'ucan/graph', with: 'api:*' }],
      description: 'Generate delegation graph'
    },
    'POST /api/graph/invocation': {
      capabilities: [{ can: 'ucan/graph', with: 'api:*' }],
      description: 'Generate invocation graph'
    }
  };
  
  /**
   * Check if a delegation has the required capability
   */
  export function hasRequiredCapability(delegation, requiredCap) {
    if (!delegation.capabilities) return false;
    
    return delegation.capabilities.some(cap => {
      const canMatch = cap.can === requiredCap.can || requiredCap.can === '*';
      const withMatch = matchesResource(cap.with, requiredCap.with);
      return canMatch && withMatch;
    });
  }
  
  /**
   * Match resources with wildcard support
   */
  function matchesResource(actual, required) {
    if (required === '*') return true;
    if (actual === required) return true;
    
    // Handle wildcards like "api:*" matching "api:parse"
    if (required.endsWith(':*')) {
      const prefix = required.slice(0, -1);
      return actual.startsWith(prefix);
    }
    
    return false;
  }