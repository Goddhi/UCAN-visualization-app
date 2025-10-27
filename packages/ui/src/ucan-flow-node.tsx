"use client";

import { memo } from "react";
import { Handle, Position } from "reactflow";

export interface UCANNodeData {
  id: string;
  issuer: string;
  audience: string;
  capabilities: string[];
  expiration?: string;
}

interface UCANFlowNodeProps {
  data: UCANNodeData;
  selected?: boolean;
}

export const UCANFlowNode = memo(({ data, selected }: UCANFlowNodeProps) => {
  return (
    <div
      className={`min-w-[280px] max-w-[320px] bg-bg-secondary border-2 rounded-lg p-4 shadow-lg transition-all ${
        selected
          ? "border-accent-primary shadow-accent-primary/30"
          : "border-border hover:border-accent-primary/50"
      }`}
    >
      {/* Top Handle */}
      <Handle
        type="target"
        position={Position.Top}
        className="w-3 h-3 !bg-accent-primary border-2 border-white"
      />

      {/* Node Header */}
      <div className="flex items-start gap-2 mb-3">
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
        <div className="min-w-0 flex-1">
          <div className="text-xs text-text-tertiary">UCAN Token</div>
          <div className="text-sm font-mono text-text-primary truncate">
            {data.id.slice(0, 12)}...
          </div>
        </div>
      </div>

      {/* Node Details */}
      <div className="space-y-2">
        <div>
          <div className="text-xs text-text-tertiary mb-1">Issuer</div>
          <div className="text-xs font-mono text-text-secondary truncate">
            {data.issuer}
          </div>
        </div>
        <div>
          <div className="text-xs text-text-tertiary mb-1">Audience</div>
          <div className="text-xs font-mono text-text-secondary truncate">
            {data.audience}
          </div>
        </div>
        <div>
          <div className="text-xs text-text-tertiary mb-1">Capabilities</div>
          <div className="flex flex-wrap gap-1">
            {data.capabilities.map((cap, idx) => (
              <span
                key={idx}
                className="px-2 py-0.5 bg-accent-primary/10 text-accent-primary text-xs rounded"
              >
                {cap}
              </span>
            ))}
          </div>
        </div>
        {data.expiration && (
          <div className="text-xs text-text-tertiary">
            Expires: {new Date(data.expiration).toLocaleDateString()}
          </div>
        )}
      </div>

      {/* Bottom Handle */}
      <Handle
        type="source"
        position={Position.Bottom}
        className="w-3 h-3 !bg-accent-primary border-2 border-white"
      />
    </div>
  );
});

UCANFlowNode.displayName = "UCANFlowNode";
