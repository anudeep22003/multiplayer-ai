import { z } from "zod";

const SESSION_KEY = "conman-session";

const nonBlankString = z.string().refine((value) => value.trim().length > 0);
const sessionSchema = z.object({
  userId: nonBlankString,
  username: nonBlankString,
});

export type Session = z.infer<typeof sessionSchema>;

export function isSession(value: unknown): value is Session {
  return sessionSchema.safeParse(value).success;
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
