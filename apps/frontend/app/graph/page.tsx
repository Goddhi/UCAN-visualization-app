"use client";

import { useState } from "react";
import { Textarea } from "@repo/ui/textarea";
import { Button } from "@repo/ui/button";
import { UCANFlow } from "@repo/ui/ucan-flow";
import { UCANNodeData } from "@repo/ui/ucan-flow-node";
import { ucanApi } from "../lib/api/client";
import type { DelegationResponse } from "../lib/api/types";
import { transformDelegationToNodeData } from "../lib/utils/transformers";

export default function GraphPage() {
  const [ucanInput, setUcanInput] = useState("");
  const [ucanData, setUcanData] = useState<DelegationResponse | null>(null);
  const [selectedNode, setSelectedNode] = useState<UCANNodeData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);

  const handleViewDelegation = async () => {
    setLoading(true);
    setError(null);

    try {
      if (file) {
        const result = await ucanApi.parseDelegationFile(file);
        setUcanData(result);
      } else if (ucanInput.trim()) {
        const result = await ucanApi.parseDelegation({
          token: ucanInput.trim(),
          format: "auto",
        });
        setUcanData(result);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to parse UCAN token",
      );
      setUcanData(null);
    } finally {
      setLoading(false);
    }
  };

  const handleClear = () => {
    setUcanInput("");
    setUcanData(null);
    setSelectedNode(null);
    setError(null);
    setFile(null);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setUcanInput("");
    }
  };

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border bg-bg-secondary px-6 py-4">
        <h1 className="text-2xl font-bold text-text-primary mb-1">
          UCAN Graph Visualizer
        </h1>
        <p className="text-sm text-text-secondary">
          Paste a UCAN token to visualize its delegation chain as an interactive
          flowchart
        </p>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Input Panel */}
        <div className="w-96 border-r border-border bg-bg-secondary p-6 flex flex-col gap-4 overflow-y-auto">
          <Textarea
            label="UCAN Token"
            placeholder="Paste your UCAN token here..."
            value={ucanInput}
            onChange={(e) => setUcanInput(e.target.value)}
            rows={8}
            className="font-mono text-sm"
            disabled={!!file}
          />

          <div className="relative">
            <div className="text-sm text-text-secondary mb-2 text-center">
              or
            </div>
            <label className="block">
              <input
                type="file"
                accept=".car,.ucan,.cbor"
                onChange={handleFileChange}
                className="hidden"
              />
              <div className="p-4 border-2 border-dashed border-border rounded-lg text-center cursor-pointer hover:border-accent-primary transition-colors">
                <svg
                  className="w-8 h-8 mx-auto mb-2 text-text-tertiary"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
                <p className="text-sm text-text-secondary">
                  {file ? file.name : "Upload .car, .ucan, or .cbor file"}
                </p>
              </div>
            </label>
          </div>

          {error && (
            <div className="p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <p className="text-sm text-red-400">{error}</p>
            </div>
          )}

          <div className="flex gap-2">
            <Button
              onClick={handleViewDelegation}
              disabled={(!ucanInput.trim() && !file) || loading}
              className="flex-1"
            >
              {loading ? (
                <>
                  <svg
                    className="w-4 h-4 animate-spin"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                  Processing...
                </>
              ) : (
                <>
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                  View Delegation
                </>
              )}
            </Button>
            <Button
              variant="secondary"
              onClick={handleClear}
              disabled={loading}
            >
              Clear
            </Button>
          </div>

          {/* Selected Node Details */}
          {selectedNode && (
            <div className="mt-4 p-4 bg-bg-tertiary border border-border rounded-lg">
              <h3 className="text-sm font-semibold text-text-primary mb-3">
                Node Details
              </h3>
              <div className="space-y-3 text-xs">
                <div>
                  <div className="text-text-tertiary mb-1">Token ID</div>
                  <div className="font-mono text-text-secondary break-all">
                    {selectedNode.id}
                  </div>
                </div>
                <div>
                  <div className="text-text-tertiary mb-1">Issuer</div>
                  <div className="font-mono text-text-secondary break-all">
                    {selectedNode.issuer}
                  </div>
                </div>
                <div>
                  <div className="text-text-tertiary mb-1">Audience</div>
                  <div className="font-mono text-text-secondary break-all">
                    {selectedNode.audience}
                  </div>
                </div>
                <div>
                  <div className="text-text-tertiary mb-1">Capabilities</div>
                  <div className="space-y-1">
                    {selectedNode.capabilities.map((cap, idx) => (
                      <div
                        key={idx}
                        className="px-2 py-1 bg-accent-primary/10 text-accent-primary rounded"
                      >
                        {cap}
                      </div>
                    ))}
                  </div>
                </div>
                {selectedNode.expiration && (
                  <div>
                    <div className="text-text-tertiary mb-1">Expiration</div>
                    <div className="text-text-secondary">
                      {new Date(selectedNode.expiration).toLocaleString()}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Instructions */}
          {!ucanData && (
            <div className="mt-auto pt-4 border-t border-border">
              <h3 className="text-sm font-semibold text-text-primary mb-2">
                How to use
              </h3>
              <ul className="text-xs text-text-secondary space-y-1 list-disc list-inside">
                <li>Paste a UCAN token in the input above</li>
                <li>Click "View Delegation" to visualize</li>
                <li>Drag nodes to rearrange the flowchart</li>
                <li>Pan the canvas by dragging empty space</li>
                <li>Zoom with mouse wheel or controls</li>
                <li>Click nodes to view details</li>
              </ul>
            </div>
          )}
        </div>

        {/* Canvas Area */}
        <div className="flex-1 relative">
          {ucanData ? (
            <UCANFlow
              data={transformDelegationToNodeData(ucanData)}
              onNodeClick={setSelectedNode}
            />
          ) : (
            <div className="flex items-center justify-center h-full">
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-r from-accent-primary to-accent-secondary flex items-center justify-center opacity-50">
                  <svg
                    className="w-8 h-8 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold text-text-primary mb-2">
                  No UCAN Token Loaded
                </h3>
                <p className="text-sm text-text-secondary max-w-md">
                  Paste a UCAN token in the panel on the left and click "View
                  Delegation" to visualize the delegation chain as an
                  interactive flowchart.
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
