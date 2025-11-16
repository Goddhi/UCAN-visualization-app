"use client";

import { useRouter } from "next/navigation";
import { Button } from "@repo/ui/button";

const guideSteps = [
  {
    title: "1. Add a UCAN",
    description: "Paste a token or upload a .car/.cbor file. Only one input is needed.",
  },
  {
    title: "2. Process",
    description: "Hit Process Token. The app calls parse, validate, and graph endpoints at once.",
  },
  {
    title: "3. Explore",
    description: "Use the canvas and side panels to inspect issuers, audiences, capabilities, and edges.",
  },
];

export default function Home() {
  const router = useRouter();

  return (
    <div className="space-y-10">
      <section className="rounded-3xl border border-border bg-bg-secondary/80 p-8 md:p-12">
        <p className="text-sm uppercase tracking-[0.2em] text-text-tertiary">
          UCAN tooling
        </p>
        <h1 className="mt-4 text-4xl font-semibold text-text-primary sm:text-5xl">
          Visualize delegation chains without the noise
        </h1>
        <p className="mt-4 max-w-2xl text-base text-text-secondary">
          Paste a token or drop a CAR file. The canvas, validator, and capability view
          sit side-by-side so you can inspect chains quickly.
        </p>
        <div className="mt-8">
          <Button size="lg" onClick={() => router.push("/graph")}>
            Open the graph
          </Button>
        </div>
      </section>

      <section className="rounded-3xl border border-border bg-bg-secondary/60 p-6">
        <h2 className="text-xl font-semibold text-text-primary">How to use it</h2>
        <div className="mt-4 grid gap-4 md:grid-cols-3">
          {guideSteps.map((step) => (
            <div key={step.title} className="rounded-2xl border border-border bg-bg-primary/30 p-4">
              <p className="text-sm font-semibold text-text-primary">{step.title}</p>
              <p className="mt-2 text-sm text-text-secondary">{step.description}</p>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
