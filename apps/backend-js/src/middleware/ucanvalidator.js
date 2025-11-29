import * as Delegation from '@ucanto/core/delegation';
import { capabilityRegistry, hasRequiredCapability } from '../swagger/capabilities.js';

/**
 * UCAN Validator Middleware
 * Validates Bearer tokens against required capabilities
 */
export function createUcanValidator() {
  return async (req, res, next) => {
    // Skip validation for health check and swagger
    const publicPaths = ['/health', '/swagger', '/swagger.json', '/api/examples/tokens'];
    if (publicPaths.some(path => req.path === path || req.path.startsWith(path))) {
      return next();
    }

    // Get Authorization header
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      return res.status(401).json({
        error: 'Missing Authorization header',
        message: 'Please provide a UCAN token in Authorization: Bearer <token> header'
      });
    }

    // Extract Bearer token
    const match = authHeader.match(/^Bearer (.+)$/);
    if (!match) {
      return res.status(401).json({
        error: 'Invalid Authorization format',
        message: 'Use format: Authorization: Bearer <base64-token>'
      });
    }

    const tokenBase64 = match[1];

    try {
      // Decode UCAN token
      const tokenBytes = Buffer.from(tokenBase64, 'base64');
      const result = await Delegation.extract(tokenBytes);
      const delegation = result.ok || result;

      // Get required capabilities for this endpoint
      const endpoint = `${req.method} ${req.path}`;
      const requirements = capabilityRegistry[endpoint];

      if (!requirements) {
        // No specific requirements, allow
        req.ucan = delegation;
        return next();
      }

      // Check if delegation has required capabilities
      const hasAllCapabilities = requirements.capabilities.every(reqCap =>
        hasRequiredCapability(delegation, reqCap)
      );

      if (!hasAllCapabilities) {
        return res.status(403).json({
          error: 'Insufficient capabilities',
          message: `Required: ${JSON.stringify(requirements.capabilities)}`,
          provided: delegation.capabilities
        });
      }

      // Attach delegation to request
      req.ucan = delegation;
      next();
    } catch (error) {
      return res.status(400).json({
        error: 'Invalid UCAN token',
        message: error.message
      });
    }
  };
}