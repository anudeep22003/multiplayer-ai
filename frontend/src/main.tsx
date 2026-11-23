import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@/index.css";
import { App } from "@/App";
import { ThemeProvider } from "@/features/theme/theme-provider";

const rootElement = document.getElementById("root");
if (!rootElement) {
  throw new Error("Root element #root was not found in index.html.");
}

createRoot(rootElement).render(
  <StrictMode>
    <ThemeProvider>
      <App />
    </ThemeProvider>
  </StrictMode>,
);
