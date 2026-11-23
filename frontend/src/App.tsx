import { isAxiosError } from "axios";
import { type FormEvent, useState } from "react";
import { Button } from "@/components/ui/button";
import { login } from "@/lib/api";
import { clearSession, readSession, saveSession } from "@/lib/session";

export function App() {
  const [session, setSession] = useState(readSession);
  const [isNew, setIsNew] = useState(false);
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (loading) return;
    const trimmedUsername = username.trim();
    if (!trimmedUsername) {
      setError("Enter a username.");
      return;
    }
    setError("");
    setLoading(true);
    try {
      const result = await login(trimmedUsername);
      saveSession(result);
      setIsNew(result.isNew);
      setSession(result);
      setUsername("");
    } catch (cause) {
      if (isAxiosError(cause) && !cause.response) {
        setError(
          "Couldn’t reach the login service. Check your connection and try again.",
        );
      } else if (cause instanceof DOMException) {
        setError(
          "Couldn’t save your session. Allow browser storage and try again.",
        );
      } else {
        setError("Login failed. Check your username and try again.");
      }
    } finally {
      setLoading(false);
    }
  }

  function handleLogout() {
    try {
      clearSession();
      setSession(null);
      setError("");
    } catch {
      setError(
        "Couldn’t clear your saved session. Allow browser storage and try again.",
      );
    }
  }

  return (
    <main className="grid min-h-svh place-items-center bg-muted/40 p-5">
      <section className="w-full max-w-sm space-y-6 rounded-xl border bg-card p-7 text-card-foreground shadow-xs">
        <h1 className="text-2xl font-semibold tracking-tight">
          {session
            ? isNew
              ? "Welcome, new user!"
              : "Welcome back!"
            : "Log in to Conman"}
        </h1>
        {session ? (
          <>
            <dl className="space-y-3">
              <div>
                <dt className="text-sm text-muted-foreground">Username</dt>
                <dd className="break-all">{session.username}</dd>
              </div>
              <div>
                <dt className="text-sm text-muted-foreground">User ID</dt>
                <dd className="break-all font-mono text-sm">
                  {session.userId}
                </dd>
              </div>
            </dl>
            <Button variant="outline" onClick={handleLogout}>
              Log out
            </Button>
          </>
        ) : (
          <form className="space-y-5" onSubmit={handleSubmit}>
            <div className="space-y-2">
              <label className="block text-sm font-medium" htmlFor="username">
                Username
              </label>
              <input
                id="username"
                name="username"
                autoComplete="username"
                className="h-10 w-full rounded-md border border-input bg-background px-3 focus-visible:outline-2"
                required
                disabled={loading}
                value={username}
                onChange={(event) => {
                  setUsername(event.target.value);
                  setError("");
                }}
              />
            </div>
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Logging in…" : "Log in"}
            </Button>
          </form>
        )}
        {error && (
          <p role="alert" className="text-sm text-destructive">
            {error}
          </p>
        )}
      </section>
    </main>
  );
}
