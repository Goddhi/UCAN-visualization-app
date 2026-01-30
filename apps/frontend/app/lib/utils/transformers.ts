import type { DelegationResponse, ProofInfo } from "../api/types";
import type { UCANNodeData } from "@repo/ui/ucan-flow-node";

export function transformDelegationToNodeData(
  delegation: DelegationResponse
): UCANNodeData & { proofs?: UCANNodeData[] } {
  return {
    id: delegation.cid,
    issuer: delegation.issuer,
    audience: delegation.audience,
    capabilities: delegation.capabilities.map(
      (cap) => `${cap.with} → ${cap.can}`
    ),
    expiration: delegation.expiration,
    proofs: (delegation.proofs?.map(transformProofToNodeData).filter(Boolean) as UCANNodeData[]) || [],
  };
}

function transformProofToNodeData(proof: ProofInfo): UCANNodeData | null {
  if (!proof.issuer || !proof.audience) {
     return {
        id: proof.cid,
        issuer: "Unknown (Link Only)",
        audience: "Unknown",
        capabilities: ["Link Only"],
        expiration: undefined,
        proofs: []
     }
  }

  return {
    id: proof.cid,
    issuer: proof.issuer,
    audience: proof.audience,
    capabilities: (proof.capabilities || []).map(
      (cap) => `${cap.with} → ${cap.can}`
    ),
    expiration: proof.expiration,
    proofs: (proof.proofs?.map(transformProofToNodeData).filter(Boolean) as UCANNodeData[]) || [],
  };
}