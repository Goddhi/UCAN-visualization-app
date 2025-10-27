// UCAN API Client for backend integration

import type {
  DelegationResponse,
  ValidationResult,
  GraphResponse,
  HealthResponse,
  ParseRequest,
  ValidateRequest,
  GraphRequest,
  ErrorResponse,
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public details?: Record<string, unknown>
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json().catch(() => ({
      error: response.statusText,
      message: 'An unknown error occurred',
    }));
    throw new ApiError(
      errorData.message || errorData.error,
      response.status,
      errorData.details
    );
  }
  return response.json();
}

export const ucanApi = {
  // Health check
  async health(): Promise<HealthResponse> {
    const response = await fetch(`${API_BASE_URL}/health`);
    return handleResponse<HealthResponse>(response);
  },

  // Parse delegation from JSON
  async parseDelegation(request: ParseRequest): Promise<DelegationResponse> {
    const response = await fetch(`${API_BASE_URL}/api/parse/delegation`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    return handleResponse<DelegationResponse>(response);
  },

  // Parse delegation from file
  async parseDelegationFile(file: File): Promise<DelegationResponse> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/api/parse/delegation/file`, {
      method: 'POST',
      body: formData,
    });
    return handleResponse<DelegationResponse>(response);
  },

  // Validate delegation chain from JSON
  async validateChain(request: ValidateRequest): Promise<ValidationResult> {
    const response = await fetch(`${API_BASE_URL}/api/validate/chain`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    return handleResponse<ValidationResult>(response);
  },

  // Validate delegation chain from file
  async validateChainFile(file: File): Promise<ValidationResult> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/api/validate/chain/file`, {
      method: 'POST',
      body: formData,
    });
    return handleResponse<ValidationResult>(response);
  },

  // Generate delegation graph from JSON
  async generateGraph(request: GraphRequest): Promise<GraphResponse> {
    const response = await fetch(`${API_BASE_URL}/api/graph/delegation`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    return handleResponse<GraphResponse>(response);
  },

  // Generate delegation graph from file
  async generateGraphFile(file: File): Promise<GraphResponse> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_BASE_URL}/api/graph/delegation/file`, {
      method: 'POST',
      body: formData,
    });
    return handleResponse<GraphResponse>(response);
  },
};

export { ApiError };
export type { ErrorResponse };
