import { ModeToggle } from "@/components/mode-toggle";
import { useTheme } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";
import { env } from "@/lib/env";

interface TokenSwatch {
  readonly label: string;
  readonly token: string;
  readonly surfaceClassName: string;
}

const TOKEN_SWATCHES: readonly TokenSwatch[] = [
  {
    label: "Background",
    token: "--background",
    surfaceClassName: "bg-background text-foreground",
  },
  {
    label: "Card",
    token: "--card",
    surfaceClassName: "bg-card text-card-foreground",
  },
  {
    label: "Muted",
    token: "--muted",
    surfaceClassName: "bg-muted text-muted-foreground",
  },
  {
    label: "Primary",
    token: "--primary",
    surfaceClassName: "bg-primary text-primary-foreground",
  },
  {
    label: "Secondary",
    token: "--secondary",
    surfaceClassName: "bg-secondary text-secondary-foreground",
  },
  {
    label: "Destructive",
    token: "--destructive",
    surfaceClassName: "bg-destructive text-destructive-foreground",
  },
];

export function App() {
  const { theme, resolvedTheme } = useTheme();

  return (
    <div className="min-h-screen bg-background text-foreground">
      <header className="border-b border-border bg-background/80 backdrop-blur-xs">
        <div className="mx-auto flex h-14 max-w-4xl items-center justify-between px-6">
          <div className="flex items-center gap-3">
            <span className="text-base font-semibold tracking-tight">
              Conman
            </span>
            <span className="rounded-md border border-border bg-muted px-2 py-0.5 font-mono text-xs text-muted-foreground">
              {theme} ({resolvedTheme})
            </span>
          </div>
          <ModeToggle />
        </div>
      </header>

      <main className="mx-auto max-w-4xl space-y-6 px-6 py-10">
        <section className="space-y-4 rounded-xl border border-border bg-card p-6 text-card-foreground shadow-xs">
          <div className="space-y-1">
            <h1 className="text-xl font-semibold tracking-tight">
              Design Token & Theme Preview
            </h1>
            <p className="text-sm text-muted-foreground">
              All surfaces, typography, borders, and actions below consume
              central OKLCH tokens from{" "}
              <code className="rounded-sm bg-muted px-1.5 py-0.5 font-mono text-xs text-foreground">
                src/index.css
              </code>
              .
            </p>
          </div>

          <div className="grid gap-3 sm:grid-cols-3">
            {TOKEN_SWATCHES.map((swatch) => (
              <div
                key={swatch.token}
                className={`rounded-lg border border-border p-3 ${swatch.surfaceClassName}`}
              >
                <div className="text-sm font-medium">{swatch.label}</div>
                <div className="font-mono text-xs opacity-80">
                  {swatch.token}
                </div>
              </div>
            ))}
          </div>

          <div className="space-y-2 pt-2">
            <h2 className="text-xs font-medium tracking-wider text-muted-foreground uppercase">
              Component Variants
            </h2>
            <div className="flex flex-wrap items-center gap-2.5">
              <Button>Primary</Button>
              <Button variant="secondary">Secondary</Button>
              <Button variant="outline">Outline</Button>
              <Button variant="ghost">Ghost</Button>
              <Button variant="destructive">Destructive</Button>
            </div>
          </div>
        </section>

        <section className="rounded-xl border border-border bg-muted/50 p-4">
          <h2 className="text-xs font-medium tracking-wider text-muted-foreground uppercase">
            Typed Runtime Environment
          </h2>
          <dl className="mt-2 grid gap-2 font-mono text-xs sm:grid-cols-2">
            <div className="rounded-md border border-border bg-card px-3 py-2 text-card-foreground">
              <dt className="text-muted-foreground">VITE_API_URL</dt>
              <dd className="mt-0.5 truncate">{env.VITE_API_URL}</dd>
            </div>
            <div className="rounded-md border border-border bg-card px-3 py-2 text-card-foreground">
              <dt className="text-muted-foreground">VITE_WS_URL</dt>
              <dd className="mt-0.5 truncate">{env.VITE_WS_URL}</dd>
            </div>
          </dl>
        </section>
      </main>
    </div>
  );
}
