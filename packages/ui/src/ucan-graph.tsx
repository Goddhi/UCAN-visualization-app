"use client";

import React, { useCallback, useState } from "react";
import ReactFlow, {
  addEdge,
  applyEdgeChanges,
  applyNodeChanges,
  Background,
  BackgroundVariant,
  Controls,
  Edge,
  MarkerType,
  Node,
  NodeChange,
  EdgeChange,
} from "reactflow";
import "reactflow/dist/style.css";
import { UcanNode } from "./ucan-node";

const nodeTypes = {
  ucanNode: UcanNode,
};

export type UcanNodeData = {
  id: string;
  label: string;
  issuer?: string;
  audience?: string;
  capabilities?: string[];
};

export type UcanGraphProps = {
  nodes?: Node<UcanNodeData>[];
  edges?: Edge[];
};

export function UcanGraph({ nodes: initialNodes = [], edges: initialEdges = [] }: UcanGraphProps) {
  const [nodes, setNodes] = useState<Node<UcanNodeData>[]>(initialNodes);
  const [edges, setEdges] = useState<Edge[]>(initialEdges);

  const onNodesChange = useCallback((changes: NodeChange[]) => setNodes((nds) => applyNodeChanges(changes, nds)), []);
  const onEdgesChange = useCallback((changes: EdgeChange[]) => setEdges((eds) => applyEdgeChanges(changes, eds)), []);

  const onConnect = useCallback((connection: any) => setEdges((eds) => addEdge({
    ...connection,
    markerEnd: { type: MarkerType.Arrow, color: 'hsl(195 100% 50%)', width: 20, height: 20 },
    style: {
      stroke: 'hsl(195 100% 50%)',
      strokeWidth: 3,
      opacity: 0.8,
      filter: 'drop-shadow(0 0 6px hsl(195 100% 50% / 0.4))'
    }
  }, eds)), []);

  return (
    <div className="w-full h-full rounded-2xl overflow-hidden border border-border/30 bg-gradient-surface">
      <ReactFlow
        nodes={nodes.map(node => ({ ...node, type: 'ucanNode' }))}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.3 }}
        defaultEdgeOptions={{
          style: {
            stroke: 'hsl(195 100% 50%)',
            strokeWidth: 3,
            opacity: 0.8,
            filter: 'drop-shadow(0 0 6px hsl(195 100% 50% / 0.4))'
          },
          markerEnd: {
            type: MarkerType.Arrow,
            color: 'hsl(195 100% 50%)',
            width: 20,
            height: 20
          },
        }}
        className="bg-transparent"
        minZoom={0.2}
        maxZoom={3}
        nodeOrigin={[0.5, 0.5]}
        connectionLineStyle={{
          stroke: 'hsl(195 100% 50%)',
          strokeWidth: 2,
          opacity: 0.6
        }}
      >
        <Background
          gap={40}
          size={2}
          variant={BackgroundVariant.Dots}
          color="hsl(var(--border))"
          className="opacity-30"
        />
        <Controls
          className="bg-glass border border-border/50 rounded-xl shadow-xl"
          style={{ left: 20, bottom: 20 }}
        />
      </ReactFlow>
    </div>
  );
}

export default UcanGraph;