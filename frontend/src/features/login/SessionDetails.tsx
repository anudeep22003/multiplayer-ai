import { Button } from "@/components/ui/button";
import type { Session } from "./session";

interface SessionDetailsProps {
  session: Session;
  isNew: boolean;
  onLogout: () => void;
}

export function SessionDetails({
  session,
  isNew,
  onLogout,
}: SessionDetailsProps) {
  return (
    <>
      <h1 className="text-2xl font-semibold tracking-tight">
        {isNew ? "Welcome, new user!" : "Welcome back!"}
      </h1>

      <dl className="space-y-3">
        <div>
          <dt className="text-sm text-muted-foreground">Username</dt>
          <dd className="break-all">{session.username}</dd>
        </div>
        <div>
          <dt className="text-sm text-muted-foreground">User ID</dt>
          <dd className="break-all font-mono text-sm">{session.userId}</dd>
        </div>
      </dl>

      <Button variant="outline" onClick={onLogout}>
        Log out
      </Button>
    </>
  );
}
