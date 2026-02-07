"use client";

import { useState } from "react";
import * as Client from "@ucanto/client";
import * as Signer from "@ucanto/principal/ed25519";
import { Button } from "@repo/ui/button";
import { Textarea } from "@repo/ui/textarea";

// Helper type for our capability rows
type CapRow = {
  can: string;
  with: string;
};

export default function BuilderPage() {
  // --- State ---
  const [issuer, setIssuer] = useState<Signer.EdSigner | null>(null);
  const [audienceDid, setAudienceDid] = useState("");
  const [capabilities, setCapabilities] = useState<CapRow[]>([
    { can: "store/add", with: "did:key:z6Mk..." },
  ]);
  const [expiration, setExpiration] = useState<string>("0"); // 0 = Infinity
  const [generatedCar, setGeneratedCar] = useState<Blob | null>(null);
  const [logs, setLogs] = useState<string[]>([]);

  // --- Actions ---

  const generateIssuer = async () => {
    const key = await Signer.generate();
    setIssuer(key);
    // Auto-fill the "with" field of the first capability with this new DID
    const newCaps = [...capabilities];
    if (newCaps[0]) newCaps[0].with = key.did();
    setCapabilities(newCaps);
    addLog(`🔑 Generated new Issuer Identity: ${key.did()}`);
  };

  const handleAddCap = () => {
    setCapabilities([...capabilities, { can: "store/add", with: issuer?.did() || "" }]);
  };

  const handleRemoveCap = (index: number) => {
    setCapabilities(capabilities.filter((_, i) => i !== index));
  };

  const handleCapChange = (index: number, field: keyof CapRow, value: string) => {
    const newCaps = [...capabilities];
    newCaps[index] = { ...newCaps[index], [field]: value };
    setCapabilities(newCaps);
  };

  const addLog = (msg: string) => setLogs((prev) => [`[${new Date().toLocaleTimeString()}] ${msg}`, ...prev]);

const handleGenerate = async () => {
    if (!issuer) return addLog("Error: No Issuer generated");
    if (!audienceDid) return addLog("Error: No Audience DID provided");

    try {
      addLog("⚙️ Generating Delegation...");

      // 1. Build the Delegation
      const delegation = await Client.delegate({
        issuer: issuer,
        audience: { did: () => audienceDid } as any, // Cast to any to mock the Principal
        capabilities: capabilities.map((c) => ({
          can: c.can,
          with: c.with,
        })) as any, // <--- FIX 1: Cast capabilities to any to bypass strict tuple checks
        expiration: expiration === "0" ? Infinity : Math.floor(Date.now() / 1000) + parseInt(expiration),
      });

      // 2. Archive to CAR
      const archive = await delegation.archive();
      if (archive.error) throw new Error("Failed to archive");

      // 3. Prepare Download
      // <--- FIX 2: Cast archive.ok to any to satisfy the Blob constructor
      const blob = new Blob([archive.ok as any], { type: "application/car" });
      
      setGeneratedCar(blob);
      addLog(` Success! Generated UCAN (${blob.size} bytes)`);
    } catch (e: any) {
      addLog(`Error: ${e.message}`);
    }
  };
  const downloadFile = () => {
    if (!generatedCar) return;
    const url = URL.createObjectURL(generatedCar);
    const a = document.createElement("a");
    a.href = url;
    a.download = "delegation.car";
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <p className="text-xs uppercase tracking-[0.3em] text-text-tertiary">Developer Tools</p>
        <h1 className="mt-2 text-3xl font-semibold text-text-primary">UCAN Builder</h1>
        <p className="mt-2 text-sm text-text-secondary">
          Visually create and sign valid UCAN delegations directly in your browser.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Left Column: Form */}
        <div className="space-y-6">
          
          {/* 1. Issuer Section */}
          <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <h3 className="text-sm font-semibold text-text-primary mb-4">1. Issuer (Who is delegating?)</h3>
            {issuer ? (
              <div className="space-y-3">
                <div className="text-xs text-text-secondary font-mono break-all bg-bg-primary/50 p-3 rounded-lg border border-border">
                  {issuer.did()}
                </div>
                <Button variant="secondary" size="sm" onClick={generateIssuer}>
                  Regenerate New Key
                </Button>
              </div>
            ) : (
              <Button onClick={generateIssuer} className="w-full">
                Generate Random Identity
              </Button>
            )}
          </section>

          {/* 2. Audience Section */}
          <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <h3 className="text-sm font-semibold text-text-primary mb-4">2. Audience (Who gets the power?)</h3>
            <Textarea
              label="Audience DID"
              placeholder="did:key:..."
              value={audienceDid}
              onChange={(e) => setAudienceDid(e.target.value)}
              rows={2}
              className="font-mono text-sm"
            />
          </section>

          {/* 3. Capabilities Section */}
          <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-semibold text-text-primary">3. Capabilities</h3>
              <Button variant="ghost" size="sm" onClick={handleAddCap}>+ Add</Button>
            </div>
            
            <div className="space-y-3">
              {capabilities.map((cap, idx) => (
                <div key={idx} className="flex gap-2 items-start">
                  <div className="grid grid-cols-[1fr_1.5fr] gap-2 flex-1">
                    <input
                      type="text"
                      placeholder="Action (e.g. store/add)"
                      className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-xs text-text-primary placeholder:text-text-tertiary focus:border-accent-primary focus:outline-none"
                      value={cap.can}
                      onChange={(e) => handleCapChange(idx, "can", e.target.value)}
                    />
                    <input
                      type="text"
                      placeholder="Resource (URI or DID)"
                      className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-xs text-text-primary placeholder:text-text-tertiary focus:border-accent-primary focus:outline-none"
                      value={cap.with}
                      onChange={(e) => handleCapChange(idx, "with", e.target.value)}
                    />
                  </div>
                  <button 
                    onClick={() => handleRemoveCap(idx)}
                    className="p-2 text-text-tertiary hover:text-error transition-colors"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          </section>

          {/* 4. Settings */}
           <section className="rounded-2xl border border-border bg-bg-secondary p-5">
            <h3 className="text-sm font-semibold text-text-primary mb-4">4. Settings</h3>
             <div className="grid grid-cols-2 gap-4">
                <div>
                   <label className="text-xs text-text-tertiary mb-1 block">Expiration (Seconds)</label>
                   <input
                      type="number"
                      placeholder="0 for Infinity"
                      className="w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-xs text-text-primary focus:border-accent-primary focus:outline-none"
                      value={expiration}
                      onChange={(e) => setExpiration(e.target.value)}
                    />
                    <span className="text-[10px] text-text-tertiary">0 = Never expires</span>
                </div>
             </div>
           </section>

          <Button 
            onClick={handleGenerate} 
            disabled={!issuer || !audienceDid}
            className="w-full py-6 text-lg"
          >
            Generate .car File
          </Button>
        </div>

        {/* Right Column: Logs & Output */}
        <div className="space-y-6">
          <section className="rounded-2xl border border-dashed border-border bg-bg-primary/20 p-5 h-full flex flex-col">
            <h3 className="text-sm font-semibold text-text-primary mb-4">Output Log</h3>
            
            <div className="flex-1 rounded-xl bg-black/5 p-4 font-mono text-[11px] overflow-y-auto max-h-[400px]">
              {logs.length === 0 && <span className="text-text-tertiary">Waiting for actions...</span>}
              {logs.map((log, i) => (
                <div key={i} className="mb-1 text-text-secondary border-b border-border/10 pb-1 last:border-0">
                  {log}
                </div>
              ))}
            </div>

            {generatedCar && (
              <div className="mt-4 pt-4 border-t border-border">
                <div className="flex items-center justify-between">
                  <span className="text-xs font-semibold text-success">Ready to download</span>
                  <span className="text-xs text-text-tertiary">{(generatedCar.size / 1024).toFixed(2)} KB</span>
                </div>
                <Button onClick={downloadFile} className="w-full mt-3" variant="secondary">
                  Download delegation.car
                </Button>
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}