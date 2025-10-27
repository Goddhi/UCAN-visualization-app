# UI Components Library

Reusable React components for the UCAN Visualization Tool, built with Tailwind CSS.

## Installation

These components are available through the monorepo workspace:

```tsx
import { Button } from "@repo/ui/button";
import { Input } from "@repo/ui/input";
import { Textarea } from "@repo/ui/textarea";
import { Sidebar, SidebarItem } from "@repo/ui/sidebar";
import { UCANTreeNode } from "@repo/ui/ucan-tree-node";
import { UCANCanvas } from "@repo/ui/ucan-canvas";
import { IconHome, IconGraph, IconInspector } from "@repo/ui/icons";
```

## Components

### Button

Multi-variant button component with different sizes and styles.

```tsx
<Button>Primary Button</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="danger">Danger</Button>

<Button size="sm">Small</Button>
<Button size="md">Medium (default)</Button>
<Button size="lg">Large</Button>

<Button disabled>Disabled</Button>
```

**Props:**
- `variant`: "primary" | "secondary" | "ghost" | "danger"
- `size`: "sm" | "md" | "lg"
- Extends all HTMLButtonElement attributes

### Input

Text input with label and error state support.

```tsx
<Input
  label="Username"
  placeholder="Enter username"
  value={value}
  onChange={(e) => setValue(e.target.value)}
/>

<Input
  label="Email"
  type="email"
  error="Invalid email address"
/>
```

**Props:**
- `label`: Optional label text
- `error`: Error message to display
- Extends all HTMLInputElement attributes

### Textarea

Multi-line text input with label and error state.

```tsx
<Textarea
  label="Description"
  placeholder="Enter description..."
  rows={5}
  value={value}
  onChange={(e) => setValue(e.target.value)}
/>

<Textarea
  label="Comments"
  error="Required field"
/>
```

**Props:**
- `label`: Optional label text
- `error`: Error message to display
- Extends all HTMLTextAreaElement attributes

### Sidebar

Collapsible navigation sidebar with retractable functionality.

```tsx
<Sidebar>
  <SidebarItem
    icon={<IconHome />}
    label="Home"
    href="/"
    active={pathname === "/"}
  />
  <SidebarItem
    icon={<IconGraph />}
    label="Graph"
    href="/graph"
    active={pathname === "/graph"}
  />
</Sidebar>
```

**Sidebar Props:**
- `children`: SidebarItem components

**SidebarItem Props:**
- `icon`: React node (typically an icon)
- `label`: Text label
- `href`: Optional link destination (creates anchor tag)
- `onClick`: Optional click handler (creates button)
- `active`: Boolean to indicate active state

### UCANTreeNode

Interactive node component for visualizing UCAN delegation chains.

```tsx
const nodeData = {
  id: "bafyreib2rxk3rybk6hj4av4xfxr5fq",
  issuer: "did:key:z6MkhaXgBZ...",
  audience: "did:key:z6MkffDZCk...",
  capabilities: ["store/add", "upload/add"],
  expiration: "2025-12-31T23:59:59Z",
  proofs: [/* nested UCAN nodes */],
};

<UCANTreeNode
  node={nodeData}
  onNodeClick={(node) => console.log("Clicked:", node)}
/>
```

**Props:**
- `node`: Object containing UCAN data
  - `id`: Token identifier
  - `issuer`: DID of issuer
  - `audience`: DID of audience
  - `capabilities`: Array of capability strings
  - `expiration`: ISO date string
  - `proofs`: Array of nested UCANNodeData objects
- `onNodeClick`: Optional callback when node is clicked

### UCANCanvas

Pan and zoom canvas for displaying UCAN graphs.

```tsx
<UCANCanvas>
  <UCANTreeNode node={ucanData} />
</UCANCanvas>
```

**Features:**
- Pan by clicking and dragging
- Zoom with mouse wheel (50% - 200%)
- Reset button to return to default view
- Scale indicator
- Pan/Zoom instructions overlay

**Props:**
- `children`: React nodes to render in the canvas

### Icons

Pre-built SVG icon components.

```tsx
<IconHome />
<IconGraph />
<IconInspector />
<IconValidate />
<IconSettings />
```

**Props:**
- `className`: Optional CSS classes (defaults to "w-5 h-5")

## Styling

All components use Tailwind CSS with a custom dark duotone color scheme:

### Color Tokens
- `bg-primary`, `bg-secondary`, `bg-tertiary`, `bg-hover`
- `accent-primary`, `accent-secondary`
- `text-primary`, `text-secondary`, `text-tertiary`
- `border`, `border-accent`

### Utilities
- `sidebar-width`: 280px
- `sidebar-collapsed-width`: 64px
- Transition durations: `fast`, `normal`, `slow`
- Custom shadows: `shadow-glow`, `shadow-glow-sm`

## Best Practices

### 1. Use Type Safety
All components are fully typed with TypeScript. Let the types guide you:

```tsx
import type { ButtonHTMLAttributes } from "react";

// Button extends ButtonHTMLAttributes
<Button onClick={handleClick} disabled={isLoading}>
  Submit
</Button>
```

### 2. Compose Components
Build complex UIs by composing smaller components:

```tsx
<Sidebar>
  <nav>
    <SidebarItem icon={<IconHome />} label="Home" href="/" />
    <SidebarItem icon={<IconGraph />} label="Graph" href="/graph" />
  </nav>
  <div className="mt-auto">
    <SidebarItem icon={<IconSettings />} label="Settings" href="/settings" />
  </div>
</Sidebar>
```

### 3. Extend with Tailwind
All components accept `className` prop for additional styling:

```tsx
<Button className="mt-4 w-full">
  Full Width Button
</Button>
```

### 4. Handle States
Components support common states like disabled, error, and active:

```tsx
<Input
  label="Email"
  value={email}
  onChange={(e) => setEmail(e.target.value)}
  error={emailError}
  disabled={isSubmitting}
/>
```

## Examples

### Form with Validation
```tsx
const [email, setEmail] = useState("");
const [error, setError] = useState("");

const handleSubmit = () => {
  if (!email.includes("@")) {
    setError("Invalid email");
    return;
  }
  // Submit logic
};

return (
  <div>
    <Input
      label="Email Address"
      type="email"
      value={email}
      onChange={(e) => setEmail(e.target.value)}
      error={error}
    />
    <Button onClick={handleSubmit} className="mt-4">
      Submit
    </Button>
  </div>
);
```

### Navigation with Active States
```tsx
"use client";

import { usePathname } from "next/navigation";

export function Navigation() {
  const pathname = usePathname();
  
  return (
    <Sidebar>
      <SidebarItem
        icon={<IconHome />}
        label="Home"
        href="/"
        active={pathname === "/"}
      />
      <SidebarItem
        icon={<IconGraph />}
        label="Graph"
        href="/graph"
        active={pathname === "/graph"}
      />
    </Sidebar>
  );
}
```

### Interactive UCAN Visualization
```tsx
const [selectedNode, setSelectedNode] = useState(null);

return (
  <div className="flex gap-4">
    <UCANCanvas>
      <UCANTreeNode
        node={ucanData}
        onNodeClick={setSelectedNode}
      />
    </UCANCanvas>
    
    {selectedNode && (
      <div className="w-80 p-4 bg-bg-secondary">
        <h3>Node Details</h3>
        <pre>{JSON.stringify(selectedNode, null, 2)}</pre>
      </div>
    )}
  </div>
);
```

## Contributing

When adding new components:

1. Create a new `.tsx` file in `/packages/ui/src`
2. Use TypeScript and proper typing
3. Follow the existing component patterns
4. Use Tailwind CSS for styling
5. Make it responsive and accessible
6. Document props and usage
7. Export from package.json if needed

## Accessibility

Components follow accessibility best practices:
- Semantic HTML elements
- ARIA labels where appropriate
- Keyboard navigation support
- Focus states
- Screen reader friendly

## Performance

- All interactive components use "use client" directive
- Optimized re-renders with React best practices
- Lazy loading where appropriate
- Memoization for expensive computations
