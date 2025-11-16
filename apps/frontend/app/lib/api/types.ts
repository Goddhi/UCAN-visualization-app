// API Types matching the backend models

export interface CapabilityInfo {
  with: string;
  can: string;
  nb?: Record<string, unknown>;
}

export interface SignatureInfo {
  algorithm?: string;
  value?: string;
}

export interface ProofInfo {
  cid: string;
  issuer: string;
  audience: string;
  capabilities: CapabilityInfo[];
  expiration?: string;
  notBefore?: string;
  proofs?: ProofInfo[];
}

export interface DelegationResponse {
  issuer: string;
  audience: string;
  subject?: string;
  capabilities: CapabilityInfo[];
  proofs: ProofInfo[];
  expiration?: string;
  notBefore?: string;
  facts?: unknown[];
  nonce?: string;
  signature?: SignatureInfo;
  cid: string;
}

export interface ValidationIssue {
  type: string;
  message: string;
  severity: 'error' | 'warning' | 'info';
  context?: Record<string, unknown>;
}

export interface ChainLink {
  level: number;
  cid: string;
  issuer: string;
  audience: string;
  capability: CapabilityInfo;
  expiration?: string;
  notBefore?: string;
  valid: boolean;
  issues?: ValidationIssue[];
}

export interface ValidationSummary {
  totalLinks?: number;
  validLinks?: number;
  invalidLinks?: number;
  warnings?: number;
  errors?: number;
}

export interface ValidationError {
  message: string;
  code?: string;
  details?: Record<string, unknown>;
}

export interface ValidationResult {
  valid: boolean;
  chain: ChainLink[];
  rootCause?: ValidationError;
  summary?: ValidationSummary;
}

export interface GraphNode {
  id: string;
  label: string;
  type: 'root' | 'intermediate' | 'leaf';
  metadata?: Record<string, unknown>;
}

export interface GraphEdge {
  source: string;
  target: string;
  capability: CapabilityInfo;
  valid: boolean;
  label: string;
}

export interface GraphResponse {
  nodes: GraphNode[];
  edges: GraphEdge[];
}

export interface ErrorResponse {
  error: string;
  message: string;
  details?: Record<string, unknown>;
  timestamp?: string;
  requestId?: string;
}

export interface ParseRequest {
  token: string;
  format?: 'base64' | 'hex' | 'auto';
}

export interface ValidateRequest {
  token: string;
  format?: 'base64' | 'hex' | 'auto';
}

export interface GraphRequest {
  token: string;
  format?: 'base64' | 'hex' | 'auto';
}

export interface HealthResponse {
  status: string;
  time: string;
  service: string;
  version: string;
}
