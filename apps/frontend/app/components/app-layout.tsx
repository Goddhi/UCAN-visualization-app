"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { IconGraph, IconHome } from "@repo/ui/icons";

const navItems = [
  { href: "/", label: "Home", icon: IconHome },
  { href: "/graph", label: "Graph", icon: IconGraph },
];

export function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-bg-primary text-text-primary">
      <header className="border-b border-border bg-bg-secondary/70 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-4">
          <Link href="/" className="flex items-center gap-2 font-semibold">
            <span className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-border text-lg">
              U
            </span>
            <span className="text-sm uppercase tracking-[0.3em] text-text-tertiary">
              UCAN
            </span>
          </Link>
          <nav className="flex items-center gap-1 text-sm">
            {navItems.map((item) => {
              const active =
                pathname === item.href ||
                (item.href !== "/" && pathname.startsWith(item.href));
              const Icon = item.icon;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`flex items-center gap-2 rounded-full border px-3 py-1.5 transition ${
                    active
                      ? "border-accent-primary bg-accent-primary/10 text-accent-primary"
                      : "border-transparent text-text-secondary hover:border-border hover:text-text-primary"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>
      </header>

      <main className="mx-auto w-full max-w-7xl px-4 py-8">{children}</main>
    </div>
  );
}
