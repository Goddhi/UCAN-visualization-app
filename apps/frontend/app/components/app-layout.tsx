"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { IconHome, IconGraph, IconBuilder } from "@repo/ui/icons"; 

const navItems = [
  { href: "/", label: "Home", icon: IconHome },
  { href: "/graph", label: "Graph", icon: IconGraph },
  { href: "/builder", label: "Builder", icon: IconBuilder },
  { href: "/inspector", label: "Invocation Inspector", icon: IconBuilder }, // <--- ADDED
];

export function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <div className="flex min-h-screen bg-bg-primary text-text-primary">
      <aside className="fixed left-0 top-0 h-screen w-64 border-r border-border bg-bg-secondary p-6">
        <div className="mb-10 flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-accent-primary text-white font-bold">U</div>
          <span className="text-lg font-bold tracking-wider">UCAN</span>
        </div>
        <nav className="space-y-2">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link key={item.href} href={item.href} className={`flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors ${isActive ? "bg-accent-primary/10 text-accent-primary" : "text-text-secondary hover:bg-bg-primary hover:text-text-primary"}`}>
                <item.icon className="h-5 w-5" />
                {item.label}
              </Link>
            );
          })}
        </nav>
      </aside>
      <main className="ml-64 flex-1 p-8"><div className="mx-auto max-w-7xl">{children}</div></main>
    </div>
  );
}