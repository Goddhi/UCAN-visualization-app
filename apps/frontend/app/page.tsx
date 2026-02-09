"use client";

import Link from "next/link";
import { Button } from "@repo/ui/button";

export default function LandingPage() {
  return (
    <div className="flex flex-col min-h-[calc(100vh-80px)]">
      {/* --- HERO SECTION --- */}
      <section className="flex flex-col items-center justify-center py-24 text-center space-y-8">
        <div className="space-y-4 max-w-3xl">
          <p className="text-sm uppercase tracking-[0.2em] text-accent-primary font-semibold">
            The Complete Toolkit
          </p>
          <h1 className="text-4xl md:text-6xl font-bold tracking-tight text-text-primary">
            Master your <span className="text-accent-primary">UCAN</span> Infrastructure
          </h1>
          <p className="text-lg text-text-secondary max-w-2xl mx-auto leading-relaxed">
            Visualize, debug, build and verify User Controlled Authorization Networks. 
          </p>
        </div>

        <div className="flex items-center gap-4 mt-4">
          {/* FIXED: Points to /graph and renamed to Visualizer */}
          <Link href="/graph">
            <Button size="lg" className="h-12 px-8 text-base bg-accent-primary hover:bg-accent-primary/90 shadow-lg shadow-accent-primary/20">
              <span className="mr-2">🔍</span> Open Visualizer
            </Button>
          </Link>
          <Link href="/builder">
            <Button variant="secondary" size="lg" className="h-12 px-8 text-base border border-border bg-bg-secondary hover:bg-bg-tertiary">
              <span className="mr-2">🛠️</span> Start Building
            </Button>
          </Link>
        </div>
      </section>

      {/* --- FEATURES GRID --- */}
      <section className="container mx-auto px-4 py-12 max-w-6xl">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-semibold text-text-primary">Powerful Features</h2>
          <p className="text-text-tertiary mt-3 text-lg">Everything you need to ship secure delegation chains.</p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          
          {/* Feature 1: Explicitly mentions Transport Blocks */}
          <FeatureCard 
            icon={<EyeIcon />}
            title="Deep Inspection"
            description="Paste raw transport blocks (tokens) to instantly decode intents. Inspect structure, verify CIDs, and resolve DIDs with one click."
          />

          {/* Feature 2: Explicitly mentions Cryptographic Verification */}
          <FeatureCard 
            icon={<ShieldCheckIcon />}
            title="Signature Verification"
            description="Cryptographically verify token signatures against issuer DIDs. Detect tampered payloads or invalid keys instantly."
          />

          {/* Feature 3: Explicitly mentions .car files */}
          <FeatureCard 
            icon={<GitBranchIcon />}
            title="Graph Visualization"
            description="Upload .car files to visualize the full delegation tree. Trace authority from root to leaf and spot broken links in the chain."
          />

          {/* Feature 4: Explicitly mentions the Visual Timeline Bar */}
          <FeatureCard 
            icon={<CalendarIcon />}
            title="Timeline Analysis"
            description="Real-time validity checks. Instantly spot expired tokens or 'not-before' violations with a visual progress bar."
          />

          {/* Feature 5: Explicitly mentions Visual/Raw Toggle */}
          <FeatureCard 
            icon={<BugIcon />}
            title="Visual Debugger"
            description="Toggle between 'Visual' and 'Raw JSON' modes. View detailed validation warnings and capability mismatches directly."
          />

          {/* Feature 6: Explicitly mentions PNG & Code Snippets */}
          <FeatureCard 
            icon={<DownloadIcon />}
            title="Export & Snippets"
            description="Download high-res PNGs of your graph. Generate ready-to-use TypeScript code snippets for invoking capabilities."
          />

        </div>
      </section>
    </div>
  );
}

// --- Components ---

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <div className="p-8 rounded-2xl border border-border bg-bg-secondary/50 hover:bg-bg-secondary hover:border-accent-primary/50 transition-all duration-300 group">
      <div className="w-12 h-12 flex items-center justify-center rounded-xl bg-accent-primary/10 text-accent-primary mb-6 group-hover:scale-110 group-hover:bg-accent-primary group-hover:text-white transition-all duration-300">
        {icon}
      </div>
      <h3 className="text-lg font-bold text-text-primary mb-3">{title}</h3>
      <p className="text-sm text-text-secondary leading-relaxed">
        {description}
      </p>
    </div>
  );
}

// --- Icons (Inline SVGs) ---

function EyeIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>;
}

function ShieldCheckIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><path d="m9 12 2 2 4-4"/></svg>;
}

function GitBranchIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="6" y1="3" x2="6" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/></svg>;
}

function CalendarIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>;
}

function BugIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect width="8" height="14" x="8" y="6" rx="4"/><path d="m19 7-3 2"/><path d="m5 7 3 2"/><path d="m19 19-3-2"/><path d="m5 19 3-2"/><path d="M20 13h-4"/><path d="M4 13h4"/><path d="m10 4 1 2"/></svg>;
}

function DownloadIcon() {
  return <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>;
}