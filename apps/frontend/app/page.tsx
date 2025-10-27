"use client";

import { useRouter } from "next/navigation";
import { Button } from "@repo/ui/button";

export default function Home() {
  const router = useRouter();

  return (
    <div className="h-full overflow-y-auto">
      {/* Hero Section */}
      <div className="relative overflow-hidden bg-gradient-to-br from-bg-primary via-bg-secondary to-bg-primary border-b border-border">
        <div className="absolute inset-0 opacity-10">
          <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-accent-primary rounded-full blur-3xl" />
          <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-accent-secondary rounded-full blur-3xl" />
        </div>

        <div className="relative max-w-7xl mx-auto px-6 py-24">
          <div className="text-center">
            <div className="inline-flex items-center gap-3 mb-6">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-r from-accent-primary to-accent-secondary flex items-center justify-center shadow-lg shadow-accent-primary/30">
                <span className="text-white font-bold text-2xl">U</span>
              </div>
            </div>

            <h1 className="text-5xl font-bold text-text-primary mb-6">
              UCAN Visualization Tool
            </h1>

            <p className="text-xl text-text-secondary max-w-2xl mx-auto mb-8">
              Parse, validate, and visualize UCAN delegation chains with ease.
              Debug decentralized authorization with visual feedback.
            </p>

            <Button onClick={() => router.push("/graph")} size="lg">
              <svg
                className="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
              Get Started
            </Button>
          </div>
        </div>
      </div>

      {/* Features */}
      <div className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold text-text-primary mb-12 text-center">
          Features
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[
            {
              icon: (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
                />
              ),
              title: "Delegation Chain Visualizer",
              description:
                "Visual tree showing who delegated what to whom with interactive exploration of trust relationships.",
            },
            {
              icon: (
                <>
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
                </>
              ),
              title: "Interactive Canvas",
              description:
                "Pan, zoom, and explore UCAN trees with an intuitive canvas inspired by Miro and Excalidraw.",
            },
            {
              icon: (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                />
              ),
              title: "Capability Breakdown",
              description:
                "Parse and display resources, abilities, and attenuations in human-readable format.",
            },
            {
              icon: (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              ),
              title: "Proof Chain Validator",
              description:
                "Check if a UCAN chain is valid with cryptographic signature verification.",
            },
          ].map((feature, idx) => (
            <div
              key={idx}
              className="p-6 bg-bg-secondary border border-border rounded-xl hover:border-accent-primary transition-all"
            >
              <div className="w-12 h-12 rounded-lg bg-gradient-to-r from-accent-primary to-accent-secondary flex items-center justify-center text-white mb-4">
                <svg
                  className="w-6 h-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  {feature.icon}
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-text-primary mb-2">
                {feature.title}
              </h3>
              <p className="text-sm text-text-secondary">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
