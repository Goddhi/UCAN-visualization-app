"use client";

import {
  useEffect,
  useMemo,
  useState,
  type ChangeEvent,
  type ReactNode,
} from "react";
import { Textarea } from "@repo/ui/textarea";
import { Button } from "@repo/ui/button";
import { UCANFlow } from "@repo/ui/ucan-flow";
import type { UCANNodeData } from "@repo/ui/ucan-flow-node";
import { ApiError, ucanApi } from "../lib/api/client";
import type {
  DelegationResponse,
  GraphEdge,
  GraphResponse,
  ValidationIssue,
  ValidationResult,
} from "../lib/api/types";
import { transformDelegationToNodeData } from "../lib/utils/transformers";

type BackendStatus = "checking" | "online" | "offline";

const backendStatusConfig: Record<
  BackendStatus,
  { label: string; dotClass: string; textClass: string }
> = {
  checking: {
    label: "Checking backend…",
    dotClass: "bg-border",
    textClass: "text-text-secondary",
  },
  online: {
    label: "Backend reachable",
    dotClass: "bg-success",
    textClass: "text-success",
  },
  offline: {
    label: "Backend offline",
    dotClass: "bg-error",
    textClass: "text-error",
  },
};

const workflowTips = [
  "Paste a UCAN token or upload a CAR/CBOR file",
  "Processing calls parse, validate, and graph endpoints",
  "Drag around the canvas and click nodes for details",
  "Zoom with your trackpad, mouse wheel, or touchpad pinch",
  "Warnings and timing issues surface in the validation panel",
];

const formatDateTime = (value?: string) =>
  value ? new Date(value).toLocaleString() : "—";

const formatDid = (value?: string) =>
  value ? `${value.slice(0, 18)}…${value.slice(-6)}` : "—";

const extractErrorMessage = (error: unknown) => {
  if (error instanceof ApiError) return error.message;
  if (error instanceof Error) return error.message;
  return "Request failed";
};

export default function GraphPage() {
  const [ucanInput, setUcanInput] = useState("");
  const [ucanData, setUcanData] = useState<DelegationResponse | null>(null);
  const [validationResult, setValidationResult] =
    useState<ValidationResult | null>(null);
  const [graphData, setGraphData] = useState<GraphResponse | null>(null);
  const [selectedNode, setSelectedNode] = useState<UCANNodeData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [graphError, setGraphError] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [backendStatus, setBackendStatus] =
    useState<BackendStatus>("checking");

  useEffect(() => {
    let mounted = true;
    ucanApi
      .health()
      .then(() => mounted && setBackendStatus("online"))
      .catch(() => mounted && setBackendStatus("offline"));
    return () => {
      mounted = false;
    };
  }, []);

  const handleViewDelegation = async () => {
    if (!file && !ucanInput.trim()) return;

    setLoading(true);
    setError(null);
    setValidationError(null);
    setGraphError(null);
    setSelectedNode(null);

    const payload = {
      token: ucanInput.trim(),
      format: undefined
    };

    try {
      const parsePromise = file
        ? ucanApi.parseDelegationFile(file)
        : ucanApi.parseDelegation(payload);
      const validatePromise = file
        ? ucanApi.validateChainFile(file)
        : ucanApi.validateChain(payload);
      const graphPromise = file
        ? ucanApi.generateGraphFile(file)
        : ucanApi.generateGraph(payload);

      const [parsed, validation, graph] = await Promise.allSettled([
        parsePromise,
        validatePromise,
        graphPromise,
      ]);

      if (parsed.status === "fulfilled") {
        setUcanData(parsed.value);
      } else {
        throw parsed.reason;
      }

      if (validation.status === "fulfilled") {
        setValidationResult(validation.value);
      } else {
        setValidationResult(null);
        setValidationError(extractErrorMessage(validation.reason));
      }

      if (graph.status === "fulfilled") {
        setGraphData(graph.value);
      } else {
        setGraphData(null);
        setGraphError(extractErrorMessage(graph.reason));
      }
    } catch (err) {
      setUcanData(null);
      setValidationResult(null);
      setGraphData(null);
      setError(
        err instanceof Error ? err.message : "Failed to process UCAN token",
      );
    } finally {
      setLoading(false);
    }
  };

  const handleClear = () => {
    setUcanInput("");
    setUcanData(null);
    setValidationResult(null);
    setGraphData(null);
    setSelectedNode(null);
    setError(null);
    setValidationError(null);
    setGraphError(null);
    setFile(null);
  };

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setUcanInput("");
    }
  };

  const flowData = useMemo(() => {
    if (!ucanData) return null;

    const root = transformDelegationToNodeData(ucanData);

    if (validationResult?.chain) {

      const hydrate = (node: any): any => {
        const updatedNode = { ...node };

        const proofDetails = validationResult.chain.find(
          (link) => link.cid === node.id
        );

        if (proofDetails) {
          updatedNode.issuer = proofDetails.issuer;
          updatedNode.audience = proofDetails.audience;
          updatedNode.capabilities = [
            `${proofDetails.capability.with} : ${proofDetails.capability.can}`
          ];
        }

        if (updatedNode.proofs && updatedNode.proofs.length > 0) {
          updatedNode.proofs = updatedNode.proofs.map((child: any) => hydrate(child));
        }

        return updatedNode;
      };

      return hydrate(root);
    }

    return root;
  }, [ucanData, validationResult]);

  const activeNodeDetails = useMemo(() => {
    // This line satisfies the linter because it USES 'selectedNode'
    if (!selectedNode) return null;

    if (validationResult?.chain) {
      const proofDetails = validationResult.chain.find(
        (link) => link.cid === selectedNode.id
      );

      if (proofDetails) {
        return {
          ...selectedNode,
          issuer: proofDetails.issuer,
          audience: proofDetails.audience,
          capabilities: [
            `${proofDetails.capability.with} : ${proofDetails.capability.can}`
          ],
          expiration: selectedNode.expiration
        };
      }
    }

    return selectedNode;
  }, [selectedNode, validationResult]);
  
  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <p className="text-xs uppercase tracking-[0.3em] text-text-tertiary">
            Graph workspace
          </p>
          <h1 className="mt-2 text-3xl font-semibold text-text-primary">
            UCAN graph visualizer
          </h1>
          <p className="mt-2 max-w-2xl text-sm text-text-secondary">
            Keep parsing, validation, and visualization together. The simpler
            layout keeps every control visible so you can focus on the data.
          </p>
        </div>
        <div className="inline-flex items-center gap-2 rounded-full border border-border px-3 py-2 text-xs text-text-secondary">
          <span
            className={`h-2 w-2 rounded-full ${backendStatusConfig[backendStatus].dotClass}`}
          />
          <span className={backendStatusConfig[backendStatus].textClass}>
            {backendStatusConfig[backendStatus].label}
          </span>
        </div>
      </div>

      <div className="grid items-start gap-6 lg:grid-cols-[minmax(0,320px)_1fr]">
        <div className="space-y-4">
          <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <Textarea
              label="UCAN Token"
              placeholder="Paste a UCAN token…"
              value={ucanInput}
              onChange={(e) => setUcanInput(e.target.value)}
              rows={8}
              className="font-mono text-sm"
              disabled={!!file}
            />
            <div className="mt-2 flex items-center justify-between text-xs text-text-tertiary">
              <span>
                {file
                  ? "Text input disabled while a file is selected"
                  : "Direct paste keeps everything local"}
              </span>
              <span>{ucanInput.length} chars</span>
            </div>
          </section>

          <section className="rounded-2xl border border-dashed border-border p-5">
            <p className="text-sm font-semibold text-text-primary">
              Upload a file instead
            </p>
            <label className="mt-3 flex cursor-pointer flex-col items-center justify-center rounded-xl border border-border bg-bg-primary/40 px-4 py-6 text-center text-sm text-text-secondary hover:border-accent-primary">
              <input
                type="file"
                accept=".car,.ucan,.cbor"
                onChange={handleFileChange}
                className="hidden"
              />
              <svg
                className="h-8 w-8 text-text-tertiary"
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
              <span className="mt-3">
                {file ? file.name : "Drop .car, .ucan, or .cbor"}
              </span>
              <span className="mt-1 text-xs text-text-tertiary">
                Max 10MB · processed locally before hitting the backend
              </span>
            </label>
            {file && (
              <Button
                variant="ghost"
                size="sm"
                className="mt-3 w-full"
                onClick={() => setFile(null)}
                disabled={loading}
              >
                Remove file
              </Button>
            )}
          </section>

          {error && (
            <div className="rounded-2xl border border-error/40 bg-error/10 px-4 py-3 text-sm text-error">
              {error}
            </div>
          )}

          <div className="flex flex-wrap gap-2">
            <Button
              onClick={handleViewDelegation}
              disabled={(!ucanInput.trim() && !file) || loading}
              className="flex-1"
            >
              {loading ? (
                <>
                  <svg
                    className="h-4 w-4 animate-spin"
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
                  Processing…
                </>
              ) : (
                <>Process token</>
              )}
            </Button>
            <Button
              variant="secondary"
              onClick={handleClear}
              disabled={loading && !!ucanData}
            >
              Reset
            </Button>
          </div>

          <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <p className="text-sm font-semibold text-text-primary">
              Workflow tips
            </p>
            <ul className="mt-3 space-y-2 text-xs text-text-secondary">
              {workflowTips.map((tip) => (
                <li key={tip} className="flex items-start gap-2">
                  <span className="mt-1 h-1.5 w-1.5 rounded-full bg-accent-primary" />
                  <span>{tip}</span>
                </li>
              ))}
            </ul>
          </section>
        </div>

        <div className="space-y-4">
          {ucanData ? (
            <>
              <section className="rounded-2xl border border-border bg-bg-secondary p-5">
                <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
                  {summaryItems.map((item) => (
                    <SummaryCard
                      key={item.label}
                      label={item.label}
                      value={item.value}
                      helper={item.helper}
                    />
                  ))}
                </div>
                <div className="mt-4 flex flex-wrap gap-2">
                  {validationResult && (
                    <StatusPill
                      tone={validationResult.valid ? "success" : "error"}
                      label={
                        validationResult.valid
                          ? "Delegation chain valid"
                          : "Delegation chain invalid"
                      }
                      helper={`links: ${
                        validationResult.summary?.totalLinks ?? 0
                      }`}
                    />
                  )}
                  {validationError && (
                    <StatusPill
                      tone="warning"
                      label="Validation unavailable"
                      helper={validationError}
                    />
                  )}
                  {graphStats && (
                    <StatusPill
                      tone="info"
                      label={`${graphStats.edges} edges`}
                      helper={`${graphData?.nodes.length ?? 0} nodes`}
                    />
                  )}
                </div>
              </section>

              <section className="rounded-2xl border border-border bg-bg-secondary p-5">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <div>
                    <h2 className="text-lg font-semibold text-text-primary">
                      Delegation tree v2
                    </h2>
                    <p className="text-xs text-text-tertiary">
                      Drag to explore or zoom to focus on a branch.
                    </p>
                  </div>
                  <span className="text-[11px] text-text-tertiary">
                    Canvas mirrors backend graph
                  </span>
                </div>
                <div className="mt-4 h-[480px] rounded-xl border border-dashed border-border bg-bg-primary/20">
                  {flowData ? (
                    <UCANFlow data={flowData} onNodeClick={setSelectedNode} />
                  ) : (
                    <div className="flex h-full items-center justify-center text-xs text-text-tertiary">
                      Unable to render graph
                    </div>
                  )}
                </div>
              </section>

              <div className="grid gap-4 lg:grid-cols-2">
                <section className="rounded-2xl border border-border bg-bg-secondary p-5">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-text-primary">
                      Node details
                    </h3>
                    <span className="text-[11px] text-text-tertiary">
                      {activeNodeDetails ? "Active" : "Click a node"}
                    </span>
                  </div>
                  {activeNodeDetails ? (
                    <div className="mt-3 space-y-3 text-xs">
                      <DetailField label="Token ID">
                        <span className="font-mono break-all text-text-secondary">
                          {activeNodeDetails.id}
                        </span>
                      </DetailField>
                      <DetailField label="Issuer">
                        <span className="font-mono break-all text-text-secondary">
                          {activeNodeDetails.issuer || "Unknown (Link Only)"}
                        </span>
                      </DetailField>
                      <DetailField label="Audience">
                        <span className="font-mono break-all text-text-secondary">
                          {activeNodeDetails.audience || "Unknown"}
                        </span>
                      </DetailField>
                      <DetailField label="Capabilities">
                        <div className="flex flex-wrap gap-2">
                          {activeNodeDetails.capabilities.length > 0 ? (
                            activeNodeDetails.capabilities.map((cap) => (
                              <span
                                key={cap}
                                className="rounded-full bg-accent-primary/10 px-2 py-0.5 text-[11px] text-accent-primary"
                              >
                                {cap}
                              </span>
                            ))
                          ) : (
                            <span className="text-text-tertiary italic">Link Only</span>
                          )}
                        </div>
                      </DetailField>
                      {activeNodeDetails.expiration && (
                        <DetailField label="Expiration">
                          {formatDateTime(activeNodeDetails.expiration)}
                        </DetailField>
                      )}
                    </div>
                  ) : (
                    <p className="mt-3 text-xs text-text-tertiary">
                      Click any node in the canvas to inspect issuer, audience,
                      and capability list.
                    </p>
                  )}
                </section>

                <section className="rounded-2xl border border-border bg-bg-secondary p-5">
                  <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-text-primary">
                      Validation
                    </h3>
                    {validationResult && (
                      <span
                        className={`text-[11px] ${
                          validationResult.valid ? "text-success" : "text-error"
                        }`}
                      >
                        {validationResult.valid ? "Valid" : "Invalid"}
                      </span>
                    )}
                  </div>
                  {validationError && (
                    <div className="mt-3 rounded-xl border border-warning/40 bg-warning/10 px-3 py-2 text-xs text-warning">
                      {validationError}
                    </div>
                  )}
                  {validationResult ? (
                    <>
                      <div className="mt-3 grid grid-cols-2 gap-3 text-xs">
                        <Stat
                          label="Links"
                          value={validationResult.summary?.totalLinks ?? 0}
                        />
                        <Stat
                          label="Warnings"
                          value={validationResult.summary?.warningCount ?? 0}
                        />
                        <Stat
                          label="Valid"
                          value={validationResult.summary?.validLinks ?? 0}
                        />
                        <Stat
                          label="Invalid"
                          value={validationResult.summary?.invalidLinks ?? 0}
                        />
                      </div>
                      <div className="mt-4 max-h-60 space-y-3 overflow-y-auto pr-1 text-xs">
                        {validationResult.chain.map((link) => (
                          <div
                            key={link.cid}
                            className="rounded-xl border border-border bg-bg-tertiary/40 p-3"
                          >
                            <div className="flex items-center justify-between text-[11px] font-semibold">
                              <span>Level {link.level}</span>
                              <span
                                className={link.valid ? "text-success" : "text-error"}
                              >
                                {link.valid ? "valid" : "invalid"}
                              </span>
                            </div>
                            <div className="mt-1 text-[11px] text-text-tertiary">
                              {formatDid(link.issuer)} → {formatDid(link.audience)}
                            </div>
                            <div className="text-[11px] text-text-secondary">
                              {link.capability.can} on {link.capability.with}
                            </div>
                            {link.issues && link.issues.length > 0 && (
                              <ValidationIssues issues={link.issues} />
                            )}
                          </div>
                        ))}
                      </div>
                      {validationResult.rootCause && (
                        <div className="mt-3 text-[11px] text-warning">
                          Root cause: {validationResult.rootCause.message}
                        </div>
                      )}
                    </>
                  ) : (
                    !validationError && (
                      <p className="mt-3 text-xs text-text-tertiary">
                        Run the validator to inspect each link, capability, and
                        time bound.
                      </p>
                    )
                  )}
                </section>
              </div>

              <section className="rounded-2xl border border-border bg-bg-secondary p-5">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-text-primary">
                    Graph data
                  </h3>
                  {graphData && (
                    <span className="text-[11px] text-text-tertiary">
                      {graphData.nodes.length} nodes
                    </span>
                  )}
                </div>
                {graphError && (
                  <div className="mt-3 rounded-xl border border-warning/40 bg-warning/10 px-3 py-2 text-xs text-warning">
                    {graphError}
                  </div>
                )}
                {graphData ? (
                  <>
                    {graphStats && (
                      <div className="mt-4 grid grid-cols-2 gap-3 text-xs">
                        <Stat label="Roots" value={graphStats.root} />
                        <Stat label="Leaves" value={graphStats.leaves} />
                        <Stat label="Intermediate" value={graphStats.intermediates} />
                        <Stat label="Edges" value={graphStats.edges} />
                      </div>
                    )}
                    <div className="mt-4 max-h-56 space-y-2 overflow-y-auto pr-1">
                      {graphData.edges.map((edge, idx) => (
                        <GraphEdgeRow
                          key={`${edge.source}-${edge.target}-${idx}`}
                          edge={edge}
                        />
                      ))}
                    </div>
                  </>
                ) : (
                  !graphError && (
                    <p className="mt-3 text-xs text-text-tertiary">
                      Graph endpoint data mirrors what the canvas renders. Run a
                      token to view the raw nodes and edges.
                    </p>
                  )
                )}
              </section>
            </>
          ) : (
            <EmptyState />
          )}
        </div>
      </div>
    </div>
  );
}

function SummaryCard({
  label,
  value,
  helper,
}: {
  label: string;
  value: ReactNode;
  helper?: ReactNode;
}) {
  return (
    <div className="rounded-xl border border-border bg-bg-primary/30 p-4">
      <div className="text-[11px] uppercase tracking-[0.2em] text-text-tertiary">
        {label}
      </div>
      <div className="mt-1 text-lg font-semibold text-text-primary">{value}</div>
      {helper && (
        <div className="mt-1 break-all text-xs text-text-secondary">{helper}</div>
      )}
    </div>
  );
}

function StatusPill({
  label,
  helper,
  tone,
}: {
  label: string;
  helper?: string;
  tone: "info" | "success" | "warning" | "error";
}) {
  const tones = {
    info: {
      border: "border-accent-primary/30",
      bg: "bg-accent-primary/10",
      text: "text-accent-primary",
      dot: "bg-accent-primary",
    },
    success: {
      border: "border-success/40",
      bg: "bg-success/10",
      text: "text-success",
      dot: "bg-success",
    },
    warning: {
      border: "border-warning/40",
      bg: "bg-warning/10",
      text: "text-warning",
      dot: "bg-warning",
    },
    error: {
      border: "border-error/40",
      bg: "bg-error/10",
      text: "text-error",
      dot: "bg-error",
    },
  } as const;

  const toneStyles = tones[tone];

  return (
    <span
      className={`inline-flex items-center gap-2 rounded-full border ${toneStyles.border} ${toneStyles.bg} ${toneStyles.text} px-3 py-1 text-xs font-medium`}
    >
      <span className={`h-2 w-2 rounded-full ${toneStyles.dot}`} />
      {label}
      {helper && (
        <span className="truncate text-text-secondary">{helper}</span>
      )}
    </span>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border border-border bg-bg-tertiary/60 p-3">
      <div className="text-[11px] text-text-tertiary uppercase tracking-[0.2em]">
        {label}
      </div>
      <div className="text-lg font-semibold text-text-primary">{value}</div>
    </div>
  );
}

function DetailField({
  label,
  children,
}: {
  label: string;
  children: ReactNode;
}) {
  return (
    <div>
      <div className="text-text-tertiary">{label}</div>
      <div className="text-text-secondary text-xs">{children}</div>
    </div>
  );
}

function ValidationIssues({ issues }: { issues: ValidationIssue[] }) {
  return (
    <div className="mt-2 space-y-1">
      {issues.map((issue, idx) => (
        <div
          key={`${issue.type}-${idx}`}
          className="flex flex-wrap gap-1 text-[11px]"
        >
          <span
            className={`font-semibold ${
              issue.severity === "warning"
                ? "text-warning"
                : issue.severity === "error"
                  ? "text-error"
                  : "text-accent-primary"
            }`}
          >
            {issue.severity ?? issue.type}
          </span>
          <span className="text-text-secondary">{issue.message}</span>
        </div>
      ))}
    </div>
  );
}

function GraphEdgeRow({ edge }: { edge: GraphEdge }) {
  const label =
    edge.label ||
    `${edge.capability.can.toUpperCase()} • ${edge.capability.with}`;
  return (
    <div className="rounded-xl border border-border bg-bg-tertiary/40 p-3">
      <div className="truncate text-xs font-semibold text-text-primary">
        {label}
      </div>
      <div className="mt-1 text-[11px] text-text-tertiary">
        {formatDid(edge.source)} → {formatDid(edge.target)}
      </div>
      <div className="text-[11px] text-text-secondary">
        {edge.capability.can} on {edge.capability.with}
      </div>
    </div>
  );
}

function EmptyState() {
  return (
    <div className="rounded-3xl border border-dashed border-border bg-bg-secondary/30 p-10 text-center">
      <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full border border-border text-2xl text-text-secondary">
        U
      </div>
      <h3 className="mt-4 text-xl font-semibold text-text-primary">
        Waiting for input
      </h3>
      <p className="mt-2 text-sm text-text-secondary">
        Paste a UCAN token or upload a file to load the graph, run validation,
        and inspect the node details side-by-side.
      </p>
    </div>
  );
}
