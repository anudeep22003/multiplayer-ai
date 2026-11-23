# Frontend Guidelines

## 1. Component Architecture

- **Keep components dumb.** Components should primarily render UI and wire up event handlers. Do not put business logic, data transformation, or API calls directly inside JSX components.
- **Extract logic into custom hooks.** If a component has non-trivial state transitions, side effects, or derived calculations, move that logic into a custom hook (`useX`).
- **Prefer composition over prop drilling.** Pass children or slots rather than threading props through 3+ layers of components.
- **Keep components small and focused.** If a component grows past ~150 lines or handles multiple distinct visual sections, split it into smaller sub-components.
- **Handle all UI states explicitly.** Any component displaying async data must cleanly handle loading, error, and empty states—prefer early returns over deeply nested ternaries.

## 2. State Management & Hooks

- **Keep state as local as possible.** Pure UI state (open/closed modals, hover, active tab, controlled input before submit) belongs in the component.
- **Lift or extract shared/domain state.** State that coordinates multiple siblings, syncs with the URL, or talks to the backend belongs in a custom hook or store—not scattered across components.
- **Derive state; don't sync it.** Never use `useEffect` to copy props into state or compute values that can be derived directly during render.
- **Minimize `useEffect`.** Only use `useEffect` to synchronize with external systems (subscriptions, DOM APIs, WebSockets, timers). Do not use it for event-driven logic—put that in event handlers.
- **Colocation rule for hooks:**
  - Used by **one** component/feature → colocate it next to that component (e.g., `ChatInput.tsx` + `useChatInput.ts`).
  - Used by **two or more** features → move it to a shared `hooks/` directory.

## 3. Data Layer & API Boundaries

- **No raw `fetch` in components.** Keep all HTTP/WebSocket calls in a dedicated `api/` or `services/` layer.
- **Validate and type at the boundary.** Type API request/response payloads explicitly. Transform backend shapes into clean UI-friendly shapes in the data layer or hook, not inside JSX.
- **Keep server state separate from client state.** Use a dedicated query/cache layer (or a clean hook abstraction) for server data rather than manually managing `isLoading`/`error`/`data` in global stores.

## 4. File Structure & Colocation

- **Organize by feature, not by file type.** Group related components, hooks, types, and helpers together inside a feature folder. Promote to shared `components/`, `hooks/`, or `lib/` only when actually reused.
- **Rule of Three.** Do not abstract a component or utility on first use. Wait until the pattern appears 2–3 times so the abstraction matches real usage.
- **Avoid barrel files (`index.ts` re-exports)** inside feature folders unless exposing a strict public API—they obscure imports and can hurt tree-shaking and bundler performance.

## 5. TypeScript & Code Hygiene

- **No `any`.** Use `unknown` and narrow it, or define explicit types.
- **Type props cleanly.** Prefer `type` or `interface` directly above the component. Avoid `React.FC`; type props directly on the function signature.
- **Use discriminated unions for state machines.** Instead of multiple booleans (`isLoading`, `isError`, `isSuccess`), prefer a union type (`{ status: 'loading' } | { status: 'error'; error: Error } | { status: 'success'; data: T }`).
- **Keep event handlers readable.** Prefix internal handlers with `handle*` (`handleSubmit`) and callback props with `on*` (`onSubmit`).

## 6. Styling & Accessibility

- **Use semantic HTML first.** Use `<button>` for actions, `<a>` for navigation, `<form>` for submissions. Never attach `onClick` to a `<div>` when a button is appropriate.
- **Treat `index.css` as the single source of truth for theming.** Use the semantic tokens defined in `index.css` (`bg-background`, `text-foreground`, `bg-muted`, `border-border`, `text-destructive`, etc.) instead of hardcoded Tailwind colors (`bg-gray-100`, `text-red-500`) or arbitrary hex/OKLCH values so light/dark modes and theme updates stay consistent.
- **Do not scatter custom CSS.** If a new reusable color, radius, or surface token is genuinely needed, define it in `index.css` (mapped to the theme) rather than creating extra `.css` files or repeating arbitrary Tailwind values. Do not modify existing theme variables in `index.css` unless explicitly asked.
- **Keep interactive elements accessible.** Icon-only buttons must have an `aria-label`. Inputs must have associated labels.

## 7. shadcn/ui Conventions

- **Keep `components/ui/` pure.** Treat `components/ui/` as dumb design-system primitives. Never add business logic, API calls, or app-specific state inside it—compose and wrap primitives in feature components instead.
- **Check primitives before building custom UI.** Reuse existing `components/ui/` components (or add them via the shadcn CLI) before hand-rolling modals, dropdowns, popovers, or form controls.
- **Merge classes with `cn()`.** Always use the `cn()` utility for conditional or overridable `className` props. Use `cva` when adding reusable variants instead of ad-hoc inline overrides or `!important`.
- **Use `asChild` for polymorphic elements.** When a `Button` or trigger needs to render as a link (`<a>` or router `<Link>`), use `asChild` rather than nesting interactive elements.
- **Preserve Radix accessibility.** Do not strip out required sub-components (e.g., `DialogTitle`, `DialogDescription`—use `sr-only` if visually hidden) so screen readers and keyboard navigation keep working.
