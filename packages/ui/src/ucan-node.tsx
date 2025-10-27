"use client";

import { useState } from "react";

export interface UCANNodeData {
  id: string;
  issuer: string;
  audience: string;
  capabilities: string[];
  expiration?: string;
  proofs?: UCANNodeData[];
}

interface UCANNodeProps {
  node: UCANNodeData;
  onNodeClick?: (node: UCANNodeData) => void;
}

export const UCANNode = ({ node, onNodeClick }: UCANNodeProps) => {
  const [isExpanded, setIsExpanded] = useState(true);
  const hasProofs = node.proofs && node.proofs.length > 0;

  return (
    <div className="flex flex-col items-center">
      {/* Node Card */}
      <div
        onClick={() => onNodeClick?.(node)}
        className="relative group min-w-[280px] max-w-[320px] bg-bg-secondary border border-border rounded-lg p-4 hover:border-accent-primary hover:shadow-lg hover:shadow-accent-primary/20 transition-all cursor-pointer"
      >
        {/* Node Header */}
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-full bg-gradient-to-r from-accent-primary to-accent-secondary flex items-center justify-center flex-shrink-0">
              <svg
                className="w-4 h-4 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                />
              </svg>
            </div>
            <div className="min-w-0">
              <div className="text-xs text-text-tertiary">UCAN Token</div>
              <div className="text-sm font-mono text-text-primary truncate">
                {node.id.slice(0, 12)}...
              </div>
            </div>
          </div>
          {hasProofs && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                setIsExpanded(!isExpanded);
              }}
              className="p-1 rounded hover:bg-bg-hover transition-colors flex-shrink-0"
            >
              <svg
                className={`w-4 h-4 text-text-secondary transition-transform ${
                  isExpanded ? "rotate-180" : ""
                }`}
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M19 9l-7 7-7-7"
                />
              </svg>
            </button>
          )}
        </div>

        {/* Node Details */}
        <div className="space-y-2">
          <div>
            <div className="text-xs text-text-tertiary mb-1">Issuer</div>
            <div className="text-xs font-mono text-text-secondary truncate">
              {node.issuer}
            </div>
          </div>
          <div>
            <div className="text-xs text-text-tertiary mb-1">Audience</div>
            <div className="text-xs font-mono text-text-secondary truncate">
              {node.audience}
            </div>
          </div>
          <div>
            <div className="text-xs text-text-tertiary mb-1">Capabilities</div>
            <div className="flex flex-wrap gap-1">
              {node.capabilities.map((cap, idx) => (
                <span
                  key={idx}
                  className="px-2 py-0.5 bg-accent-primary/10 text-accent-primary text-xs rounded"
                >
                  {cap}
                </span>
              ))}
            </div>
          </div>
          {node.expiration && (
            <div className="text-xs text-text-tertiary">
              Expires: {new Date(node.expiration).toLocaleDateString()}
            </div>
          )}
        </div>

        {/* Connection Point */}
        {hasProofs && (
          <div className="absolute -bottom-3 left-1/2 -translate-x-1/2 w-6 h-6 bg-bg-secondary border border-border rounded-full flex items-center justify-center">
            <div className="w-2 h-2 bg-accent-primary rounded-full" />
          </div>
        )}
      </div>

      {/* Child Nodes */}
      {hasProofs && isExpanded && (
        <div className="relative mt-8">
          {/* Connecting Line */}
          <div className="absolute top-0 left-1/2 -translate-x-1/2 w-0.5 h-4 bg-border" />

          <div className="flex gap-8 pt-4">
            {node.proofs?.map((proof) => (
              <div key={proof.id} className="relative">
                {/* Connecting Line to Parent */}
                <div className="absolute -top-4 left-1/2 -translate-x-1/2 w-0.5 h-4 bg-border" />
                <UCANNode node={proof} onNodeClick={onNodeClick} />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
