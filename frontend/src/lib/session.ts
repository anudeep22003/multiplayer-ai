const SESSION_KEY = "conman-session";

export interface Session {
  userId: string;
  username: string;
}

export function isSession(value: unknown): value is Session {
  return (
    typeof value === "object" &&
    value !== null &&
    "userId" in value &&
    typeof value.userId === "string" &&
    value.userId.trim() !== "" &&
    "username" in value &&
    typeof value.username === "string" &&
    value.username.trim() !== ""
  );
}

export function readSession(): Session | null {
  try {
    const stored = localStorage.getItem(SESSION_KEY);
    if (!stored) return null;
    const session: unknown = JSON.parse(stored);
    return isSession(session) ? session : null;
  } catch {
    // Unavailable storage or invalid JSON means there is no usable session.
    return null;
  }
}

export function saveSession({ userId, username }: Session) {
  localStorage.setItem(SESSION_KEY, JSON.stringify({ userId, username }));
}

export function clearSession() {
  localStorage.removeItem(SESSION_KEY);
}
