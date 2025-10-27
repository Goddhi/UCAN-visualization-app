"use client";

import { usePathname } from "next/navigation";
import { Sidebar, SidebarItem } from "@repo/ui/sidebar";
import { IconHome, IconGraph, IconSettings } from "@repo/ui/icons";

export function AppLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  return (
    <div className="flex h-screen overflow-hidden bg-bg-primary">
      <Sidebar>
        <SidebarItem
          icon={<IconHome />}
          label="Home"
          href="/"
          active={pathname === "/"}
        />
        <SidebarItem
          icon={<IconGraph />}
          label="UCAN Graph"
          href="/graph"
          active={pathname === "/graph"}
        />
        <div className="mt-auto pt-4 border-t border-border">
          <SidebarItem
            icon={<IconSettings />}
            label="Settings"
            href="/settings"
            active={pathname === "/settings"}
          />
        </div>
      </Sidebar>
      <main className="flex-1 overflow-hidden">{children}</main>
    </div>
  );
}
