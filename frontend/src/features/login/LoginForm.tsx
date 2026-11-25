import { type SubmitEvent, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface LoginFormProps {
  loading: boolean;
  onLogin: (username: string) => void;
  onEdit: () => void;
}

export function LoginForm({ loading, onLogin, onEdit }: LoginFormProps) {
  const [username, setUsername] = useState("");

  function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    onLogin(username);
  }

  return (
    <>
      <h1 className="text-2xl font-semibold tracking-tight">
        Log in to Multiplayer AI
      </h1>

      <form className="space-y-5" onSubmit={handleSubmit}>
        <div className="space-y-2">
          <label className="block text-sm font-medium" htmlFor="username">
            Username
          </label>
          <Input
            id="username"
            name="username"
            autoComplete="username"
            className="h-10 bg-background dark:bg-background"
            required
            disabled={loading}
            value={username}
            onChange={(event) => {
              setUsername(event.target.value);
              onEdit();
            }}
          />
        </div>

        <Button type="submit" className="w-full" disabled={loading}>
          {loading ? "Logging in…" : "Log in"}
        </Button>
      </form>
    </>
  );
}
