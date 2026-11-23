import { createContext, type ReactNode, useContext } from "react";
import {
  type Theme,
  type ThemeContextValue,
  useThemePreference,
} from "./useThemePreference";

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

interface ThemeProviderProps {
  readonly children: ReactNode;
  readonly defaultTheme?: Theme;
  readonly storageKey?: string;
}

export function ThemeProvider({
  children,
  defaultTheme,
  storageKey,
}: ThemeProviderProps) {
  const value = useThemePreference(defaultTheme, storageKey);

  return (
    <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
  );
}

export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider.");
  }

  return context;
}
