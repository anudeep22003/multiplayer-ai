import { z } from "zod";

const SESSION_KEY = "conman-session";

export const SessionSchema = z.object({
  userId: z.string().trim().min(1),
  username: z.string().trim().min(1),
});

export type Session = z.infer<typeof SessionSchema>;

export function readSession(): Session | null {
  try {
    const stored = localStorage.getItem(SESSION_KEY);
    if (!stored) return null;
    const parsed = SessionSchema.safeParse(JSON.parse(stored));
    return parsed.success ? parsed.data : null;
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
