"use client";

export default function SettingsPage() {
  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="border-b border-border bg-bg-secondary px-6 py-4">
        <h1 className="text-2xl font-bold text-text-primary mb-1">Settings</h1>
        <p className="text-sm text-text-secondary">Configure your preferences</p>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        <div className="max-w-4xl mx-auto">
          <div className="bg-bg-secondary border border-border rounded-xl p-6">
            <h2 className="text-lg font-semibold text-text-primary mb-4">Appearance</h2>

            <div className="space-y-4">
              <div className="flex items-center justify-between py-3 border-b border-border">
                <div>
                  <div className="text-sm font-medium text-text-primary">Theme</div>
                  <div className="text-xs text-text-secondary">Dark duotone theme (default)</div>
                </div>
                <div className="px-3 py-1.5 bg-bg-tertiary border border-border rounded text-sm text-text-secondary">
                  Dark
                </div>
              </div>

              <div className="flex items-center justify-between py-3">
                <div>
                  <div className="text-sm font-medium text-text-primary">Accent Color</div>
                  <div className="text-xs text-text-secondary">Primary accent color</div>
                </div>
                <div className="flex gap-2">
                  <div className="w-8 h-8 rounded bg-gradient-to-r from-accent-primary to-accent-secondary border-2 border-white/20" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
