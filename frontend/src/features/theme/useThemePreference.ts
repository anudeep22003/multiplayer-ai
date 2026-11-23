import { useEffect, useMemo, useState } from "react";
export type Theme = "light" | "dark" | "system";
export type ResolvedTheme = "light" | "dark";

const DEFAULT_STORAGE_KEY = "conman-theme";
const DARK_MEDIA_QUERY = "(prefers-color-scheme: dark)";

export interface ThemeContextValue {
  readonly theme: Theme;
  readonly resolvedTheme: ResolvedTheme;
  readonly setTheme: (theme: Theme) => void;
}

export function isTheme(value: string | null): value is Theme {
  return value === "light" || value === "dark" || value === "system";
}

function getSystemTheme(): ResolvedTheme {
  if (typeof window === "undefined") {
    return "light";
  }
  return window.matchMedia(DARK_MEDIA_QUERY).matches ? "dark" : "light";
}

function readStoredTheme(storageKey: string, fallback: Theme): Theme {
  if (typeof window === "undefined") {
    return fallback;
  }
  try {
    const stored = window.localStorage.getItem(storageKey);
    return isTheme(stored) ? stored : fallback;
  } catch {
    return fallback;
  }
}

export function useThemePreference(
  defaultTheme: Theme = "system",
  storageKey = DEFAULT_STORAGE_KEY,
) {
  const [theme, setThemeState] = useState<Theme>(() =>
    readStoredTheme(storageKey, defaultTheme),
  );
  const [systemTheme, setSystemTheme] = useState<ResolvedTheme>(getSystemTheme);

  useEffect(() => {
    const mediaQuery = window.matchMedia(DARK_MEDIA_QUERY);
    const handleChange = (event: MediaQueryListEvent) => {
      setSystemTheme(event.matches ? "dark" : "light");
    };

    setSystemTheme(mediaQuery.matches ? "dark" : "light");
    mediaQuery.addEventListener("change", handleChange);
    return () => mediaQuery.removeEventListener("change", handleChange);
  }, []);

  const resolvedTheme: ResolvedTheme = theme === "system" ? systemTheme : theme;

  useEffect(() => {
    const root = window.document.documentElement;
    const isDark = resolvedTheme === "dark";
    root.classList.toggle("dark", isDark);
    root.style.colorScheme = resolvedTheme;
  }, [resolvedTheme]);

  return useMemo<ThemeContextValue>(
    () => ({
      theme,
      resolvedTheme,
      setTheme: (nextTheme: Theme) => {
        try {
          window.localStorage.setItem(storageKey, nextTheme);
        } catch {
          // Continue updating in-memory state if storage is unavailable.
        }
        setThemeState(nextTheme);
      },
    }),
    [theme, resolvedTheme, storageKey],
  );
}
