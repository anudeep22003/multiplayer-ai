import { LoginForm } from "@/features/login/LoginForm";
import { SessionDetails } from "@/features/login/SessionDetails";
import { useLogin } from "@/features/login/useLogin";

export function App() {
  const { state, handleLogin, handleLogout, clearLoginError } = useLogin();
  const error =
    state.status === "failed" || state.status === "signedIn"
      ? state.error
      : undefined;

  return (
    <main className="grid min-h-svh place-items-center bg-muted/40 p-5">
      <section className="w-full max-w-sm space-y-6 rounded-xl border bg-card p-7 text-card-foreground shadow-xs">
        {state.status === "signedIn" ? (
          <SessionDetails
            session={state.session}
            isNew={state.isNew}
            onLogout={handleLogout}
          />
        ) : (
          <LoginForm
            loading={state.status === "submitting"}
            onLogin={handleLogin}
            onEdit={clearLoginError}
          />
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
