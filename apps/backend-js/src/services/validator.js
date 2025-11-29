export class ValidatorService {
  constructor(parserService) {
    this.parser = parserService;
  }

  async validateChain(tokenBytes) {
    try {
      const chain = await this.parser.parseDelegationChain(tokenBytes);
      const chainLinks = [];
      const allIssues = [];

      for (const delegation of chain) {
        const link = this.#validateDelegation(delegation);
        chainLinks.push(link);
        allIssues.push(...link.issues);
      }

      const summary = this.#buildSummary(chainLinks);
      let rootCause = null;
      
      if (summary.invalidLinks > 0) {
        rootCause = this.#findRootCause(allIssues, chainLinks[0]);
      }

      return {
        valid: summary.invalidLinks === 0,
        chain: chainLinks,
        rootCause,
        summary
      };
    } catch (error) {
      return {
        valid: false,
        rootCause: {
          type: 'parse_error',
          message: `Failed to parse UCAN: ${error.message}`
        },
        summary: {
          totalLinks: 0,
          validLinks: 0,
          invalidLinks: 0,
          warningCount: 0
        }
      };
    }
  }

  #validateDelegation(delegation) {
    const issues = [];
    const now = new Date();

    if (delegation.expiration) {
      const expiration = new Date(delegation.expiration);
      
      if (expiration < now) {
        const timeExpired = Math.floor((now - expiration) / 60000);
        issues.push({
          type: 'expired',
          message: `UCAN expired ${this.#formatDuration(timeExpired)} ago`,
          severity: 'error'
        });
      } else if (expiration < new Date(now.getTime() + 24 * 60 * 60 * 1000)) {
        const timeUntil = Math.floor((expiration - now) / 60000);
        issues.push({
          type: 'expiring_soon',
          message: `UCAN expires in ${this.#formatDuration(timeUntil)}`,
          severity: 'warning'
        });
      }
    }

    if (delegation.notBefore) {
      const notBefore = new Date(delegation.notBefore);
      if (notBefore > now) {
        issues.push({
          type: 'not_yet_valid',
          message: `UCAN not valid until ${notBefore.toISOString()}`,
          severity: 'error'
        });
      }
    }

    if (delegation.capabilities.length === 0) {
      issues.push({
        type: 'no_capabilities',
        message: 'Delegation has no capabilities',
        severity: 'warning'
      });
    }

    const capability = delegation.capabilities[0] || {
      with: '', can: '', nb: {}
    };

    const valid = this.#countErrors(issues) === 0;

    return {
      level: delegation.level,
      cid: delegation.cid,
      issuer: delegation.issuer,
      audience: delegation.audience,
      capability,
      expiration: delegation.expiration,
      notBefore: delegation.notBefore,
      valid,
      issues
    };
  }

  #buildSummary(links) {
    const summary = {
      totalLinks: links.length,
      validLinks: 0,
      invalidLinks: 0,
      warningCount: 0
    };

    for (const link of links) {
      if (link.valid) {
        summary.validLinks++;
      } else {
        summary.invalidLinks++;
      }

      for (const issue of link.issues) {
        if (issue.severity === 'warning') {
          summary.warningCount++;
        }
      }
    }

    return summary;
  }

  #findRootCause(issues, firstLink) {
    for (const issue of issues) {
      if (issue.severity === 'error') {
        return {
          type: issue.type,
          message: issue.message,
          link: {
            issuer: firstLink.issuer,
            audience: firstLink.audience
          }
        };
      }
    }
    return null;
  }

  #countErrors(issues) {
    return issues.filter(i => i.severity === 'error').length;
  }

  #formatDuration(minutes) {
    if (minutes < 60) return `${minutes}m`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours < 24) return `${hours}h ${mins}m`;
    const days = Math.floor(hours / 24);
    const hrs = hours % 24;
    return `${days}d ${hrs}h`;
  }
}
