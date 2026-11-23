type EnvKey = "VITE_API_URL" | "VITE_WS_URL";

function requireUrlEnv(
  key: EnvKey,
  allowedProtocols: readonly string[],
): string {
  const raw = import.meta.env[key]?.trim();
  if (!raw) {
    throw new Error(
      `Missing required environment variable: ${key}. Copy .env.example to .env.local and configure it.`,
    );
  }

  let parsed: URL;
  try {
    parsed = new URL(raw);
  } catch {
    throw new Error(`Invalid URL in environment variable ${key}: "${raw}".`);
  }

  if (!allowedProtocols.includes(parsed.protocol)) {
    throw new Error(
      `Invalid protocol "${parsed.protocol}" for ${key}: expected ${allowedProtocols.join(" or ")}.`,
    );
  }

  return raw;
}

export const VITE_API_URL = requireUrlEnv("VITE_API_URL", ["http:", "https:"]);
export const VITE_WS_URL = requireUrlEnv("VITE_WS_URL", ["ws:", "wss:"]);

export const env = {
  VITE_API_URL,
  VITE_WS_URL,
} as const;
