"use client";

import { ReactNode, useRef, useState, useEffect } from "react";

interface UCANCanvasProps {
  children: ReactNode;
}

export const UCANCanvas = ({ children }: UCANCanvasProps) => {
  const canvasRef = useRef<HTMLDivElement>(null);
  const [isPanning, setIsPanning] = useState(false);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [startPos, setStartPos] = useState({ x: 0, y: 0 });
  const [scale, setScale] = useState(1);

  const handleMouseDown = (e: React.MouseEvent) => {
    if (e.button === 0) { // Left mouse button
      setIsPanning(true);
      setStartPos({
        x: e.clientX - position.x,
        y: e.clientY - position.y,
      });
    }
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (isPanning) {
      setPosition({
        x: e.clientX - startPos.x,
        y: e.clientY - startPos.y,
      });
    }
  };

  const handleMouseUp = () => {
    setIsPanning(false);
  };

  const handleWheel = (e: React.WheelEvent) => {
    e.preventDefault();
    const delta = e.deltaY * -0.001;
    const newScale = Math.min(Math.max(0.5, scale + delta), 2);
    setScale(newScale);
  };

  const handleReset = () => {
    setPosition({ x: 0, y: 0 });
    setScale(1);
  };

  return (
    <div className="relative w-full h-full overflow-hidden">
      {/* Canvas Controls */}
      <div className="absolute top-4 right-4 z-10 flex gap-2">
        <button
          onClick={() => setScale(Math.min(scale + 0.1, 2))}
          className="p-2 bg-card border border-border rounded-lg hover:bg-secondary transition-colors"
          title="Zoom In"
        >
          <svg className="w-5 h-5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
        </button>
        <button
          onClick={() => setScale(Math.max(scale - 0.1, 0.5))}
          className="p-2 bg-card border border-border rounded-lg hover:bg-secondary transition-colors"
          title="Zoom Out"
        >
          <svg className="w-5 h-5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
          </svg>
        </button>
        <button
          onClick={handleReset}
          className="p-2 bg-card border border-border rounded-lg hover:bg-secondary transition-colors"
          title="Reset View"
        >
          <svg className="w-5 h-5 text-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>

      {/* Scale Indicator */}
      <div className="absolute bottom-4 right-4 z-10 px-3 py-1.5 bg-card border border-border rounded-lg">
        <span className="text-xs text-muted-foreground">{Math.round(scale * 100)}%</span>
      </div>

      {/* Pan/Zoom Instructions */}
      <div className="absolute bottom-4 left-4 z-10 px-3 py-1.5 bg-card border border-border rounded-lg">
        <span className="text-xs text-muted-foreground">
          Drag to pan • Scroll to zoom
        </span>
      </div>

      {/* Canvas */}
      <div
        ref={canvasRef}
        className={`w-full h-full ${isPanning ? "cursor-grabbing" : "cursor-grab"}`}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onWheel={handleWheel}
      >
        <div
          style={{
            transform: `translate(${position.x}px, ${position.y}px) scale(${scale})`,
            transformOrigin: "0 0",
            transition: isPanning ? "none" : "transform 0.1s ease-out",
          }}
          className="w-full h-full flex items-start justify-center pt-20"
        >
          {children}
        </div>
      </div>
    </div>
  );
};
