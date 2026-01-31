"use client";

import { useCallback, useEffect, useMemo } from "react"; 
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  MarkerType,
  BackgroundVariant,
} from "reactflow";
import "reactflow/dist/style.css";
import { UCANFlowNode, UCANNodeData } from "./ucan-flow-node";

interface UCANFlowProps {
  data: UCANNodeData & { proofs?: UCANNodeData[] };
  onNodeClick?: (node: UCANNodeData) => void;
}

const nodeTypes = {
  ucanNode: UCANFlowNode,
};

// ... (Keep buildFlowData function exactly as is) ...
function buildFlowData(
  ucanData: UCANNodeData & { proofs?: UCANNodeData[] },
  parentId: string | null = null,
  xOffset: number = 0,
  yOffset: number = 0,
  level: number = 0
): { nodes: Node[]; edges: Edge[] } {
  const nodes: Node[] = [];
  const edges: Edge[] = [];

  nodes.push({
    id: ucanData.id,
    type: "ucanNode",
    position: { x: xOffset, y: yOffset },
    data: ucanData,
  });

  if (parentId) {
    edges.push({
      id: `${parentId}-${ucanData.id}`,
      source: parentId,
      target: ucanData.id,
      type: "smoothstep",
      animated: true,
      style: { stroke: "#6366f1", strokeWidth: 2 },
      markerEnd: { type: MarkerType.ArrowClosed, color: "#6366f1" },
    });
  }

  if (ucanData.proofs && ucanData.proofs.length > 0) {
    const childCount = ucanData.proofs.length;
    const horizontalSpacing = 400;
    const verticalSpacing = 250;
    const totalWidth = (childCount - 1) * horizontalSpacing;
    const startX = xOffset - totalWidth / 2;

    ucanData.proofs.forEach((proof, index) => {
      const childX = startX + index * horizontalSpacing;
      const childY = yOffset + verticalSpacing;
      const childData = buildFlowData(proof, ucanData.id, childX, childY, level + 1);
      nodes.push(...childData.nodes);
      edges.push(...childData.edges);
    });
  }

  return { nodes, edges };
}

export const UCANFlow = ({ data, onNodeClick }: UCANFlowProps) => {
  const { nodes: newNodes, edges: newEdges } = useMemo(
    () => buildFlowData(data, null, 400, 50),
    [data]
  );

  const [nodes, setNodes, onNodesChange] = useNodesState(newNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(newEdges);

  // --- FIX START: Sync state when data changes ---
  useEffect(() => {
    setNodes(newNodes);
    setEdges(newEdges);
  }, [newNodes, newEdges, setNodes, setEdges]);
  // --- FIX END ---

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  const onNodeClickHandler = useCallback(
    (_event: React.MouseEvent, node: Node) => {
      if (onNodeClick) {
        onNodeClick(node.data as UCANNodeData);
      }
    },
    [onNodeClick]
  );

  return (
    <div className="w-full h-full bg-bg-primary">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeClick={onNodeClickHandler}
        nodeTypes={nodeTypes}
        fitView
        minZoom={0.5}
        maxZoom={1.5}
        defaultViewport={{ x: 0, y: 0, zoom: 1 }}
        className="bg-bg-primary"
      >
        <Background variant={BackgroundVariant.Dots} gap={20} size={1} color="#2a2a38" />
        <Controls className="!bg-bg-secondary !border-border" />
      </ReactFlow>
    </div>
  );
};