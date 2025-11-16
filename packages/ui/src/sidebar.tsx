"use client";

import { ReactNode, useState } from "react";

interface SidebarProps {
  children?: ReactNode;
}

interface SidebarItemProps {
  icon: ReactNode;
  label: string;
  href?: string;
  active?: boolean;
  onClick?: () => void;
}

export const Sidebar = ({ children }: SidebarProps) => {
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <>
      <aside
        className={`fixed left-0 top-0 h-full bg-bg-secondary/80 backdrop-blur-xl border-r border-border/50 transition-all duration-300 z-40 ${
          isCollapsed ? "w-20" : "w-72"
        }`}
      >
        <div className="flex flex-col h-full">
          {/* Header with Logo */}
          <div className="p-6 border-b border-border/50">
            <div
              className={`flex items-center ${isCollapsed ? "justify-center" : "gap-3"}`}
            >
              <div className="relative group">
                <div className="absolute inset-0 bg-gradient-to-r from-accent-primary to-accent-secondary rounded-xl blur-md opacity-60 group-hover:opacity-80 transition-opacity" />
                <div className="relative w-10 h-10 rounded-xl bg-gradient-to-r from-accent-primary to-accent-secondary flex items-center justify-center shadow-lg">
                  <svg
                    className="w-6 h-6 text-white"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"
                    />
                  </svg>
                </div>
              </div>
              {!isCollapsed && (
                <div className="flex flex-col">
                  <span className="font-bold text-lg text-text-primary">
                    UCAN
                  </span>
                  <span className="text-xs text-text-tertiary font-medium">
                    Visualization
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Toggle Button */}
          <button
            onClick={() => setIsCollapsed(!isCollapsed)}
            className="absolute -right-3 top-8 w-6 h-6 rounded-full bg-bg-tertiary border border-border flex items-center justify-center hover:bg-accent-primary hover:border-accent-primary hover:scale-110 transition-all duration-200 group shadow-lg"
            aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
          >
            <svg
              className={`w-3 h-3 text-text-secondary group-hover:text-white transition-all duration-300 ${
                isCollapsed ? "rotate-180" : ""
              }`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2.5}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>

          {/* Navigation */}
          <nav className="flex-1 px-3 py-6 overflow-y-auto">
            {!isCollapsed && (
              <div className="mb-4 px-3">
                <span className="text-xs font-semibold text-text-tertiary uppercase tracking-wider">
                  Navigation
                </span>
              </div>
            )}
            <div className="space-y-1">{children}</div>
          </nav>

          {/* Footer */}
          <div className="p-4 border-t border-border/50 bg-bg-tertiary/30">
            {!isCollapsed ? (
              <div className="flex items-center gap-3 px-2">
                <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-accent-primary/20 to-accent-secondary/20 border border-accent-primary/30 flex items-center justify-center">
                  <span className="text-xs font-bold text-accent-primary">
                    v1
                  </span>
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-xs font-medium text-text-primary">
                    Version 0.1.0
                  </div>
                  <div className="text-xs text-text-tertiary">Beta Release</div>
                </div>
              </div>
            ) : (
              <div className="flex justify-center">
                <div className="w-2 h-2 rounded-full bg-accent-primary animate-pulse" />
              </div>
            )}
          </div>
        </div>
      </aside>
      {/* Spacer */}
      <div
        className={`transition-all duration-300 ${isCollapsed ? "w-20" : "w-72"}`}
      />
    </>
  );
};

export const SidebarItem = ({
  icon,
  label,
  href,
  active = false,
  onClick,
}: SidebarItemProps) => {
  const Component = href ? "a" : "button";
  const props = href ? { href } : { onClick, type: "button" as const };

  return (
    <Component
      {...props}
      className={`group relative flex items-center gap-3 w-full px-3 py-3 rounded-xl transition-all duration-200 ${
        active
          ? "bg-gradient-to-r from-accent-primary to-accent-secondary text-white shadow-lg shadow-accent-primary/25"
          : "text-text-secondary hover:bg-bg-hover hover:text-text-primary"
      }`}
    >
      {/* Active Indicator */}
      {active && (
        <div className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-8 bg-white rounded-r-full" />
      )}

      {/* Icon */}
      <span
        className={`flex-shrink-0 transition-all duration-200 ${
          active
            ? "text-white"
            : "text-text-tertiary group-hover:text-accent-primary group-hover:scale-110"
        }`}
      >
        {icon}
      </span>

      {/* Label */}
      <span className="font-medium text-sm truncate">{label}</span>

      {/* Hover Effect */}
      {!active && (
        <div className="absolute inset-0 rounded-xl bg-gradient-to-r from-accent-primary/0 to-accent-secondary/0 group-hover:from-accent-primary/5 group-hover:to-accent-secondary/5 transition-all duration-200" />
      )}

      {/* Arrow Indicator on Active */}
      {active && (
        <svg
          className="ml-auto w-4 h-4 text-white opacity-70"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5l7 7-7 7"
          />
        </svg>
      )}
    </Component>
  );
};
